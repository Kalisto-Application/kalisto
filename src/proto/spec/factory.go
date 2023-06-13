package spec

import (
	"fmt"
	"kalisto/src/models"
	"kalisto/src/proto/compiler"

	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) FromRegistry(reg *compiler.Registry) (spec models.Spec, err error) {
	for _, fd := range reg.Descriptors {
		services := fd.GetServices()
		specServices := make([]models.Service, len(services))

		for i, service := range services {
			methods := service.GetMethods()
			specMethods := make([]models.Method, len(methods))

			for j, method := range methods {

				input := method.GetInputType()
				inputFields := input.GetFields()
				specInputFields := make([]models.Field, len(inputFields))

				for k, inputField := range inputFields {
					specField, err := f.newField(inputField)
					if err != nil {
						return spec, fmt.Errorf("proto spec: failed to create new field: %w", err)
					}

					specInputFields[k] = specField
				}

				specMethods[j] = models.Method{
					Name:     method.GetName(),
					FullName: method.GetFullyQualifiedName(),
					RequestMessage: models.Message{
						Name:     input.GetName(),
						FullName: input.GetFullyQualifiedName(),
						Fields:   specInputFields,
					},
					Kind: models.NewCommunicationKind(method.IsClientStreaming(), method.IsServerStreaming()),
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

	return spec, err
}

func (f *Factory) newField(fd *desc.FieldDescriptor) (_ models.Field, err error) {
	var dataType models.DataType
	var enum []string
	var fields []models.Field
	var isCollection bool
	var collectionKey *models.Field
	var oneOf []models.Field
	var defaultValue string

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
			panic("not implemented")
		}

		message := fd.GetMessageType()
		mFields := message.GetFields()
		fields = make([]models.Field, len(mFields))
		for i := range mFields {
			field, err := f.newField(mFields[i])
			if err != nil {
				return field, err
			}

			fields[i] = field
		}
		defaultValue = "{}"

	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		return models.Field{}, models.ErrProto2NotSupported
	}

	return models.Field{
		Name:          fd.GetName(),
		FullName:      fd.GetFullyQualifiedJSONName(),
		Type:          dataType,
		DefaultValue:  defaultValue,
		Enum:          enum,
		IsCollection:  isCollection,
		CollectionKey: collectionKey,
		OneOf:         oneOf,
		Fields:        fields,
	}, nil
}
