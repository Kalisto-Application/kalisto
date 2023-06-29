package spec

import (
	"fmt"
	"kalisto/src/models"
	"kalisto/src/proto/compiler"
	"strings"
	"sync"

	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Factory struct {
	mx    sync.RWMutex
	links map[string]models.Message
}

func NewFactory() *Factory {
	return &Factory{links: make(map[string]models.Message)}
}

func (f *Factory) FromRegistry(reg *compiler.Registry) (spec models.Spec, err error) {
	f.mx.Lock()
	defer f.mx.Unlock()
	if err := f.linkMessages(reg); err != nil {
		return spec, err
	}

	for _, fd := range reg.Descriptors {
		services := fd.GetServices()
		specServices := make([]models.Service, len(services))

		for i, service := range services {
			methods := service.GetMethods()
			specMethods := make([]models.Method, len(methods))

			for j, method := range methods {

				input := method.GetInputType()
				msg, ok := f.links[input.GetFullyQualifiedName()]
				if !ok {
					return spec, fmt.Errorf("type %s not found: %w", input.GetFullyQualifiedName(), err)
				}
				requestExample := f.makeRequestExample(msg)

				specMethods[j] = models.Method{
					Name:           method.GetName(),
					FullName:       method.GetFullyQualifiedName(),
					RequestMessage: msg,
					Kind:           models.NewCommunicationKind(method.IsClientStreaming(), method.IsServerStreaming()),
					RequestExample: requestExample,
				}
			}

			specServices[i] = models.Service{
				Name:     service.GetName(),
				FullName: service.GetFullyQualifiedName(),
				Methods:  specMethods,
				Package:  service.GetFile().GetPackage(),
			}
		}

		spec.Services = append(spec.Services, specServices...)
	}

	spec.Links = f.links
	return spec, err
}

func (f *Factory) newField(fd *desc.FieldDescriptor) (_ models.Field, err error) {
	var dataType models.DataType
	var enum []string
	var isCollection bool
	var collectionKey *models.Field
	var oneOf []models.Field
	var defaultValue string
	var link string

	switch fd.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		dataType = models.DataTypeBool
		defaultValue = "false"
	case descriptorpb.FieldDescriptorProto_TYPE_INT32, descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_UINT32, descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64, descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32, descriptorpb.FieldDescriptorProto_TYPE_FIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64, descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		dataType = models.DataTypeInt
		defaultValue = "0"
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT, descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		dataType = models.DataTypeFloat
		defaultValue = "0.0"
	case descriptorpb.FieldDescriptorProto_TYPE_STRING, descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		dataType = models.DataTypeString
		defaultValue = `""`
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		dataType = models.DataTypeEnum
		v := fd.GetEnumType().GetValues()
		enum = make([]string, len(v))
		for i := range v {
			enum[i] = v[i].GetName()
		}
		defaultValue = v[0].GetName()
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		defaultValue = "{}"
		message := fd.GetMessageType()
		linkKey := message.GetFullyQualifiedName()
		if _, ok := f.links[linkKey]; !ok {
			if err := f.linkMessageFields(message, linkKey); err != nil {
				return models.Field{}, err
			}
		}

		link = linkKey

	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		return models.Field{}, models.ErrProto2NotSupported
	}

	if fd.IsRepeated() {
		isCollection = true
		defaultValue = "[]"
	}

	if fd.IsMap() {
		isCollection = true
		defaultValue = "{}"
		key, err := f.newField(fd.GetMapKeyType())
		if err != nil {
			return models.Field{}, err
		}
		collectionKey = &key
	}

	if oneOf := fd.GetOneOf(); oneOf != nil {
		return models.Field{}, nil
		panic("not implemented")
	}

	specField := models.Field{
		Name:          fd.GetName(),
		FullName:      fd.GetFullyQualifiedName(),
		Type:          dataType,
		DefaultValue:  defaultValue,
		Enum:          enum,
		IsCollection:  isCollection,
		CollectionKey: collectionKey,
		OneOf:         oneOf,
		Message:       link,
	}
	return specField, nil
}

func (f *Factory) linkMessages(reg *compiler.Registry) (err error) {
	for _, fd := range reg.Descriptors {
		mTypes := fd.GetMessageTypes()
		for _, mt := range mTypes {
			fullName := mt.GetFullyQualifiedName()
			if _, ok := f.links[fullName]; ok {
				continue
			}

			if err := f.linkMessageFields(mt, fullName); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *Factory) linkMessageFields(mt *desc.MessageDescriptor, key string) error {
	mFields := mt.GetFields()
	f.links[key] = models.Message{
		Name:     mt.GetName(),
		FullName: mt.GetFullyQualifiedName(),
		Fields:   make([]models.Field, len(mFields)),
	}

	for i, mf := range mFields {
		mField, err := f.newField(mf)
		if err != nil {
			return err
		}
		f.links[key].Fields[i] = mField
	}

	return nil
}

func (f *Factory) makeRequestExample(m models.Message) string {
	var buf strings.Builder
	buf.WriteString("{\n")

	space := 2
	for _, field := range m.Fields {
		var v string
		switch field.Type {
		case models.DataTypeString:
			v = `"string"`
		}

		line := fmt.Sprintf("%s%s: %s,", strings.Repeat(" ", space), field.Name, v)
		buf.WriteString(line)
	}

	buf.WriteString("}")
	return buf.String()
}
