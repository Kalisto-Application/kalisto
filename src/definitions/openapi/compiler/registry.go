package compiler

import (
	"fmt"
	"kalisto/src/models"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

type Registry struct {
	docs []*openapi3.T

	links map[string]models.Message
	mx    sync.RWMutex
}

func NewRegistry(docs []*openapi3.T) *Registry {
	return &Registry{docs: docs, links: make(map[string]models.Message)}
}

func (r *Registry) Schema() (models.Spec, error) {
	services := make([]models.Service, 0, len(r.docs))

	for _, doc := range r.docs {
		if doc.Info == nil {
			continue
		}

		methods := make([]models.Method, 0)
		for pathStr, path := range doc.Paths {
			params := path.Parameters
			for _, op := range []*openapi3.Operation{path.Connect, path.Delete, path.Get, path.Head, path.Options, path.Patch, path.Post, path.Put, path.Trace} {
				if op == nil {
					continue
				}
				fullName := doc.Info.Title + "." + op.OperationID + "Request"
				methodFullName := fmt.Sprintf("%s.%s", pathStr, op.OperationID)

				pathParams := make([]*openapi3.Parameter, 0)
				queryParams := make([]*openapi3.Parameter, 0)
				bodyParams := make([]*openapi3.Parameter, 0)

				for _, param := range append(params, op.Parameters...) {
					switch param.Value.In {
					case "path":
						pathParams = append(pathParams, param.Value)
					case "query":
						queryParams = append(queryParams, param.Value)
					case "body":
						bodyParams = append(bodyParams, param.Value)
					}
				}

				fields := make([]models.Field, 0)
				pathFields := make([]models.Field, 0, len(pathParams))
				queryFields := make([]models.Field, 0, len(queryParams))
				bodyFields := make([]models.Field, 0, len(bodyParams))
				for _, param := range pathParams {
					f, err := r.makeFieldFromParam(param)
					if err != nil {
						return models.Spec{}, fmt.Errorf("failed to make field from path param: %w", err)
					}

					pathFields = append(pathFields, f)
				}
				for _, param := range queryParams {
					f, err := r.makeFieldFromParam(param)
					if err != nil {
						return models.Spec{}, fmt.Errorf("failed to make field from path param: %w", err)
					}

					queryFields = append(queryFields, f)
				}
				for _, param := range bodyParams {
					f, err := r.makeFieldFromParam(param)
					if err != nil {
						return models.Spec{}, fmt.Errorf("failed to make field from path param: %w", err)
					}

					bodyFields = append(bodyFields, f)
				}
				if len(pathFields) > 0 {
					pathFullName := fullName + ".path"
					pathMsg := models.Message{
						Name:     "path",
						FullName: pathFullName,
						Fields:   pathFields,
					}
					r.mx.Lock()
					r.links[pathFullName] = pathMsg
					r.mx.Unlock()

					fields = append(fields, models.Field{
						Name:     "path",
						FullName: pathFullName,
						Type:     models.DataTypeStruct,
						Message:  pathFullName,
					})
				}
				if len(queryFields) > 0 {
					queryFullName := fullName + ".query"
					queryMsg := models.Message{
						Name:     "query",
						FullName: queryFullName,
						Fields:   pathFields,
					}
					r.mx.Lock()
					r.links[queryFullName] = queryMsg
					r.mx.Unlock()

					fields = append(fields, models.Field{
						Name:     "query",
						FullName: queryFullName,
						Type:     models.DataTypeStruct,
						Message:  queryFullName,
					})
				}
				if len(bodyFields) > 0 {
					bodyFullName := fullName + ".body"
					bodyMsg := models.Message{
						Name:     "body",
						FullName: bodyFullName,
						Fields:   pathFields,
					}
					r.mx.Lock()
					r.links[bodyFullName] = bodyMsg
					r.mx.Unlock()

					fields = append(fields, models.Field{
						Name:     "body",
						FullName: bodyFullName,
						Type:     models.DataTypeStruct,
						Message:  bodyFullName,
					})
				}

				if len(bodyParams) == 0 && op.RequestBody != nil {
					for _, v := range op.RequestBody.Value.Content {
						_ = v
						panic("not implemented")
					}
				}

				requestMsg := models.Message{
					Name:     op.OperationID + "Request",
					FullName: fullName,
					Fields:   fields,
				}

				r.mx.Lock()
				r.links[requestMsg.FullName] = requestMsg
				r.mx.Unlock()

				methods = append(methods, models.Method{
					Name:           op.OperationID,
					FullName:       methodFullName,
					Kind:           models.CommunicationKindSimple,
					RequestMessage: requestMsg,
				})
			}
		}

		services = append(services, models.Service{
			Name:        doc.Info.Title,
			DisplayName: doc.Info.Title,
			Package:     doc.Info.Version,
			FullName:    fmt.Sprintf("%s:%s", doc.Info.Title, doc.Info.Version),
			Methods:     methods,
		})
	}

	return models.Spec{
		Services: services,
		Links:    r.Links(),
	}, nil
}

func (r *Registry) makeFieldFromParam(param *openapi3.Parameter) (models.Field, error) {
	if param.Schema != nil {
		return r.makeFieldFromSchema(param.Name, param.Schema)
	}

	return models.Field{
		Name: param.Name,
		Type: models.DataTypeString,
	}, nil
}

func (r *Registry) makeFieldFromSchema(name string, schemaRef *openapi3.SchemaRef) (models.Field, error) {
	switch schemaRef.Value.Type {
	case "object":
		fields := make([]models.Field, 0)

		for name, prop := range schemaRef.Value.Properties {
			field, err := r.makeFieldFromSchema(name, prop)
			if err != nil {
				return field, fmt.Errorf("failed to make field from schema: %w", err)
			}

			fields = append(fields, field)
		}

		r.mx.Lock()
		r.links[schemaRef.Ref] = models.Message{
			Name:     schemaRef.Ref,
			FullName: schemaRef.Ref,
			Fields:   fields,
		}
		r.mx.Unlock()
		return models.Field{
			Name:     name,
			FullName: name,
			Type:     models.DataTypeStruct,
			Message:  schemaRef.Ref,
		}, nil
	case "string":
		return models.Field{
			Name:     name,
			FullName: name,
			Type:     models.DataTypeString,
		}, nil
	case "boolean":
		return models.Field{
			Name:     name,
			FullName: name,
			Type:     models.DataTypeBool,
		}, nil
	case "number":
		switch schemaRef.Value.Format {
		case "float":
			return models.Field{
				Name:     name,
				FullName: name,
				Type:     models.DataTypeFloat32,
			}, nil
		default:
			return models.Field{
				Name:     name,
				FullName: name,
				Type:     models.DataTypeFloat64,
			}, nil
		}
	case "integer":
		switch schemaRef.Value.Format {
		case "int32":
			return models.Field{
				Name:     name,
				FullName: name,
				Type:     models.DataTypeInt32,
			}, nil
		default:
			return models.Field{
				Name:     name,
				FullName: name,
				Type:     models.DataTypeInt64,
			}, nil
		}

	case "array":
		return r.makeFieldFromSchema(name, schemaRef.Value.Items)
	}

	return models.Field{}, fmt.Errorf("unknown schema type: %s", schemaRef.Value.Type)
}

func (r *Registry) NewResponseMessage(methodFullName string) (*dynamic.Message, error) {
	return nil, nil
}
func (r *Registry) GetInputType(methodFullName string) (*desc.MessageDescriptor, error) {
	return nil, nil
}
func (r *Registry) GetOutputType(methodFullName string) (*desc.MessageDescriptor, error) {
	return nil, nil
}
func (r *Registry) MethodPath(methodFullName string) (string, error) {
	return "", nil
}
func (r *Registry) Links() map[string]models.Message {
	r.mx.RLock()
	defer r.mx.RUnlock()

	links := make(map[string]models.Message, len(r.links))
	for k, v := range r.links {
		links[k] = v
	}

	return links
}
