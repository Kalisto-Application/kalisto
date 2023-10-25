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
	services := make([]models.Service, len(r.docs))
	links := make(map[string]models.Message)

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
				methodFullName := fmt.Sprintf("%s.%s", pathStr, op.OperationID)

				// filter in path
				// filteri in query
				// filter in body
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
					fields = append(fields, models.Field{
						Name:     "path",
						FullName: fmt.Sprintf("%s.%s", methodFullName, "path"),
						Type:     models.DataTypeStruct,
					})
				}
				if len(queryFields) > 0 {
					fields = append(fields, models.Field{
						Name:     "query",
						FullName: fmt.Sprintf("%s.%s", methodFullName, "query"),
						Type:     models.DataTypeStruct,
					})
				}
				if len(bodyFields) > 0 {
					fields = append(fields, models.Field{
						Name:     "body",
						FullName: fmt.Sprintf("%s.%s", methodFullName, "body"),
						Type:     models.DataTypeStruct,
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
					FullName: path.Ref + op.OperationID + "Request",
					Fields:   fields,
				}
				_ = requestMsg

				methods = append(methods, models.Method{
					Name:     op.OperationID,
					FullName: methodFullName,
					Kind:     models.CommunicationKindSimple,
					// RequestMessage:  msg,s,
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
		Links:    links,
	}, nil
}

func (r *Registry) makeFieldFromParam(param *openapi3.Parameter) (models.Field, error) {
	if param.Schema != nil {
		return r.makeFieldFromSchema(param.Name, param.Schema.Ref, param.Schema.Value)
	}

	return models.Field{
		Name: param.Name,
		Type: models.DataTypeString,
	}, nil
}

func (r *Registry) makeFieldFromSchema(name string, ref string, schema *openapi3.Schema) (models.Field, error) {
	switch schema.Type {
	case "object":
		fields := make([]models.Field, 0)

		for name, prop := range schema.Properties {
			field, err := r.makeFieldFromSchema(name, prop.Ref, prop.Value)
			if err != nil {
				return field, fmt.Errorf("failed to make field from schema: %w", err)
			}

			fields = append(fields, field)
		}

		r.mx.Lock()
		r.links[ref] = models.Message{
			Name:     ref,
			FullName: ref,
			Fields:   fields,
		}
		r.mx.Unlock()
		return models.Field{
			Name:     name,
			FullName: name,
			Type:     models.DataTypeStruct,
			Message:  ref,
		}, nil
	case "string":
		return models.Field{
			Name:     name,
			FullName: name,
			Type:     models.DataTypeString,
		}, nil
	case "array":

	}

	return models.Field{}, fmt.Errorf("unknown schema type: %s", schema.Type)
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
