package spec

import (
	"fmt"
	"kalisto/src/models"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) FromRegistry(reg *protoregistry.Files) (spec models.Spec, extErr error) {
	reg.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		services := fd.Services()
		specServices := make([]models.Service, services.Len())

		for i := 0; i < services.Len(); i++ {
			service := services.Get(i)
			methods := service.Methods()
			specMethods := make([]models.Method, methods.Len())

			for j := 0; j < methods.Len(); j++ {
				method := methods.Get(j)

				input := method.Input()
				inputFields := input.Fields()
				specInputFields := make([]models.Field, inputFields.Len())

				for k := 0; k < inputFields.Len(); k++ {
					specField, err := f.newField(inputFields.Get(k))
					if err != nil {
						extErr = fmt.Errorf("proto spec: failed to create new field: %w", err)
						return false
					}

					specInputFields[k] = specField
				}

				specMethods[j] = models.Method{
					Name: string(method.Name()),
					RequestMessage: models.Message{
						Name:   string(input.Name()),
						Fields: specInputFields,
					},
					Kind: models.NewCommunicationKind(method.IsStreamingClient(), method.IsStreamingServer()),
				}
			}

			specServices[i] = models.Service{
				Name:    string(service.Name()),
				Methods: specMethods,
			}

		}

		spec.Services = append(spec.Services, specServices...)
		return true
	})

	return spec, extErr
}

func (f *Factory) newField(fd protoreflect.FieldDescriptor) (_ models.Field, err error) {
	var dataType models.DataType
	var enum []string
	var fields []models.Field
	var isCollection bool
	var collectionKey *models.Field
	var oneOf []models.Field
	var defaultValue string

	switch fd.Kind() {
	case protoreflect.BoolKind:
		dataType = models.DataTypeBool
		defaultValue = "false"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind, protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
		dataType = models.DataTypeInt
		defaultValue = "0"
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		dataType = models.DataTypeFloat
		defaultValue = "0.0"
	case protoreflect.StringKind, protoreflect.BytesKind:
		dataType = models.DataTypeString
		defaultValue = `""`
	case protoreflect.EnumKind:
		dataType = models.DataTypeEnum
		v := fd.Enum().Values()
		enum = make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			enum[i] = string(v.Get(i).Name())
		}
		defaultValue = string(v.Get(0).Name())
	case protoreflect.MessageKind:
		if fd.IsList() {
			isCollection = true
			defaultValue = "[]"
		}

		if fd.IsMap() {
			isCollection = true
			defaultValue = "{}"
			key, err := f.newField(fd.MapKey())
			if err != nil {
				return models.Field{}, err
			}
			collectionKey = &key
		}

		if oneOf := fd.ContainingOneof(); oneOf != nil {
			panic("not implemented")
		}

		message := fd.Message()
		mFields := message.Fields()
		fields = make([]models.Field, mFields.Len())
		for i := 0; i < mFields.Len(); i++ {
			field, err := f.newField(mFields.Get(i))
			if err != nil {
				return field, err
			}

			fields[i] = field
		}
		defaultValue = "{}"

	case protoreflect.GroupKind:
		return models.Field{}, models.ErrProto2NotSupported
	}

	return models.Field{
		Name:          string(fd.Name()),
		Type:          dataType,
		DefaultValue:  defaultValue,
		Enum:          enum,
		IsCollection:  isCollection,
		CollectionKey: collectionKey,
		OneOf:         oneOf,
		Fields:        fields,
	}, nil
}
