package spec

import (
	"errors"
	"fmt"
	"kalisto/src/models"
	"kalisto/src/proto/compiler"
	"strings"
	"sync"

	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Factory struct {
	mx    sync.RWMutex
	links map[string]models.Message
}

func NewFactory() *Factory {
	return &Factory{links: make(map[string]models.Message)}
}

var ErrFieldIsOneOf = errors.New("field is one of")

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
				requestExample := f.makeRequestExample(make(map[string]bool), msg, 2, "")

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

func (f *Factory) newField(fd *desc.FieldDescriptor, oneOfName string, set map[string]bool) (_ models.Field, err error) {
	var dataType models.DataType
	var enum []int32
	var mapKey *models.Field
	var mapValue *models.Field
	var oneOfs []models.Field
	var link string
	repeated := fd.IsRepeated()
	name := fd.GetName()
	fullName := fd.GetFullyQualifiedName()

	switch fd.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		dataType = models.DataTypeBool
	case descriptorpb.FieldDescriptorProto_TYPE_INT32, descriptorpb.FieldDescriptorProto_TYPE_SINT32, descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		dataType = models.DataTypeInt32
	case descriptorpb.FieldDescriptorProto_TYPE_INT64, descriptorpb.FieldDescriptorProto_TYPE_SINT64, descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		dataType = models.DataTypeInt64
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32, descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		dataType = models.DataTypeUint32
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64, descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		dataType = models.DataTypeUint64
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		dataType = models.DataTypeFloat32
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		dataType = models.DataTypeFloat64
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		dataType = models.DataTypeString
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		dataType = models.DataTypeBytes
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		dataType = models.DataTypeEnum
		v := fd.GetEnumType().GetValues()
		enum = make([]int32, len(v))
		for i := range v {
			enum[i] = v[i].GetNumber()
		}
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		switch fd.GetMessageType().GetFullyQualifiedName() {
		case "google.protobuf.Timestamp":
			dataType = models.DataTypeDate
		case "google.protobuf.Duration":
			dataType = models.DataTypeDuration
		default:
			if fd.IsMap() {
				repeated = false
				key, err := f.newField(fd.GetMapKeyType(), "", set)
				if err != nil {
					return models.Field{}, err
				}
				mapKey = &key

				valueField, err := f.newField(fd.GetMapValueType(), "", set)
				if err != nil {
					return models.Field{}, err
				}
				mapValue = &valueField
			}

			dataType = models.DataTypeStruct
			message := fd.GetMessageType()
			linkKey := message.GetFullyQualifiedName()
			if _, ok := f.links[linkKey]; !ok {
				if err := f.linkMessageFields(message, linkKey); err != nil {
					return models.Field{}, err
				}
			}

			link = linkKey
		}

	case descriptorpb.FieldDescriptorProto_TYPE_GROUP:
		return models.Field{}, models.ErrProto2NotSupported
	}

	if oneOf := fd.GetOneOf(); oneOf != nil && oneOfName == "" {
		if set[oneOf.GetName()] {
			return models.Field{}, ErrFieldIsOneOf
		}
		dataType = models.DataTypeOneOf
		set[oneOf.GetName()] = true
		oneOfs = make([]models.Field, len(oneOf.GetChoices()))

		name = oneOf.GetName()
		fullName = oneOf.GetFullyQualifiedName()

		for i, choice := range oneOf.GetChoices() {
			oneOfField, err := f.newField(choice, oneOf.GetName(), set)
			if err != nil {
				return models.Field{}, err
			}

			oneOfs[i] = oneOfField
		}
	}

	specField := models.Field{
		Name:     name,
		FullName: fullName,
		Type:     dataType,
		Enum:     enum,
		Repeated: repeated,
		MapKey:   mapKey,
		MapValue: mapValue,
		OneOf:    oneOfs,
		Message:  link,
	}
	return specField, nil
}

func (f *Factory) linkMessages(reg *compiler.Registry) (err error) {
	protoregistry.GlobalTypes.RangeMessages(func(t protoreflect.MessageType) bool {
		mt, err := desc.WrapMessage(t.Descriptor())
		if err != nil {
			fmt.Println("failed to wrap message descriptor: ", string(t.Descriptor().FullName()), err.Error())
		}

		fullName := mt.GetFullyQualifiedName()
		if _, ok := f.links[fullName]; ok {
			return true
		}

		if err := f.linkMessageFields(mt, fullName); err != nil {
			fmt.Println("failed to link default msg: ", fullName, err.Error())
			return false
		}

		return true
	})

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
	fields := make([]models.Field, 0, len(mFields))
	f.links[key] = models.Message{
		Name:     mt.GetName(),
		FullName: mt.GetFullyQualifiedName(),
	}

	set := make(map[string]bool)
	for _, mf := range mFields {
		mField, err := f.newField(mf, "", set)
		if err != nil {
			if errors.Is(err, ErrFieldIsOneOf) {
				continue
			}
			return err
		}
		fields = append(fields, mField)
	}
	f.links[key] = models.Message{
		Name:     mt.GetName(),
		FullName: mt.GetFullyQualifiedName(),
		Fields:   fields,
	}

	return nil
}

func (f *Factory) makeRequestExample(set map[string]bool, m models.Message, space int, parent string) string {
	var buf strings.Builder
	buf.WriteString("{\n")

	for _, field := range m.Fields {
		setKey := fmt.Sprintf("%s:%s:%s", m.FullName, parent, field.FullName)
		if field.Message != "" && set[setKey] {
			continue
		}
		set[setKey] = true
		v := f.makeExampleValue(set, field, space, m.FullName)
		if v == "" {
			continue
		}

		tpl := "%s%s: %s,\n"
		if field.Type == models.DataTypeOneOf {
			tpl = "%s%s: %s"
		}
		line := fmt.Sprintf(tpl, strings.Repeat(" ", space), field.Name, v)
		buf.WriteString(line)
	}

	closeBracket := fmt.Sprintf("%s}", strings.Repeat(" ", space-2))
	buf.WriteString(closeBracket)
	return buf.String()
}

func (f *Factory) makeExampleValue(set map[string]bool, field models.Field, space int, parent string) string {
	if field.Repeated {
		fieldCp := field
		fieldCp.Repeated = false
		v := f.makeExampleValue(set, fieldCp, space, "")
		return fmt.Sprintf("[%s]", v)
	}

	var v string
	switch field.Type {
	case models.DataTypeString:
		v = `"string"`
	case models.DataTypeBool:
		v = `true`
	case models.DataTypeInt32, models.DataTypeInt64, models.DataTypeUint32, models.DataTypeUint64:
		v = `1`
	case models.DataTypeFloat32, models.DataTypeFloat64:
		v = `3.14`
	case models.DataTypeBytes:
		v = `"{json: true}"`
	case models.DataTypeEnum:
		if len(field.Enum) == 0 {
			return ""
		}
		v = fmt.Sprintf(`%d`, field.Enum[0])
	case models.DataTypeDuration:
		v = "1576800000000000"
	case models.DataTypeDate:
		v = "Date.now()"
	case models.DataTypeStruct:
		if field.MapKey != nil && field.MapValue != nil {
			key := f.makeExampleValue(set, *field.MapKey, space, "")
			value := f.makeExampleValue(set, *field.MapValue, space, "")
			v = fmt.Sprintf(`{%s: %s}`, key, value)
		} else {
			link, ok := f.links[field.Message]
			if !ok {
				return ""
			}
			if strings.HasPrefix(parent, field.FullName) {
				return ""
			}
			parent += ":" + field.FullName
			v = f.makeRequestExample(set, link, space+2, parent)
		}
	case models.DataTypeOneOf:
		var oneOfBuf strings.Builder
		for i, one := range field.OneOf {
			oneV := f.makeExampleValue(set, one, space, field.FullName)
			oneV = fmt.Sprintf("{\"%s\": %s},\n", one.Name, oneV)
			if i != 0 {
				oneV = strings.Repeat(" ", space) + "// " + field.Name + ": " + oneV
			}
			oneOfBuf.WriteString(oneV)
		}
		v = oneOfBuf.String()
	}

	return v
}
