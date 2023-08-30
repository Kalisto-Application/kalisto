package spec

import (
	"errors"
	"fmt"
	"kalisto/src/models"
	"kalisto/src/proto/compiler"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

				input := method.GetInputType().GetFullyQualifiedName()
				msg, ok := f.links[input]
				if !ok {
					return spec, fmt.Errorf("type %s not found: %w", input, err)
				}
				requestExample := f.makeRequestExample(make(map[string]bool), msg, 2, "")
				out := method.GetOutputType().GetFullyQualifiedName()
				responseMsg, ok := f.links[out]
				if !ok {
					return spec, fmt.Errorf("type %s not found: %w", out, err)
				}

				specMethods[j] = models.Method{
					Name:            method.GetName(),
					FullName:        method.GetFullyQualifiedName(),
					RequestMessage:  msg,
					ResponseMessage: responseMsg,
					Kind:            models.NewCommunicationKind(method.IsClientStreaming(), method.IsServerStreaming()),
					RequestExample:  requestExample,
				}
			}

			var serviceName string
			serviceNameParts := strings.Split(service.GetFullyQualifiedName(), ".")
			if len(serviceNameParts) == 1 {
				serviceName = serviceNameParts[0]
			} else {
				serviceName = strings.Join(serviceNameParts[len(serviceNameParts)-2:], ".")
			}

			specServices[i] = models.Service{
				Name:        service.GetName(),
				DisplayName: serviceName,
				FullName:    service.GetFullyQualifiedName(),
				Methods:     specMethods,
				Package:     service.GetFile().GetPackage(),
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
			oneV = fmt.Sprintf("{%s: %s},\n", one.Name, oneV)
			if i != 0 {
				oneV = field.Name + ": " + oneV
				lines := strings.Split(oneV, "\n")
				for i, line := range lines {
					if strings.TrimSpace(line) == "" {
						continue
					}
					lines[i] = strings.Repeat(" ", space) + "// " + line
				}
				oneV = strings.Join(lines, "\n")
			}
			oneOfBuf.WriteString(oneV)
		}
		v = oneOfBuf.String()
	}

	return v
}

func (f *Factory) MessageAsJsString(spec models.Message, msg *dynamic.Message) (string, error) {
	return f.asJs(spec, msg, 2)
}

func (f *Factory) asJs(m models.Message, msg *dynamic.Message, space int) (string, error) {
	var buf strings.Builder
	buf.WriteString("{\n")

	for _, field := range m.Fields {
		var protoValue any
		var err error
		tpl := strings.Repeat(" ", space) + "%s: %s,\n"
		if len(field.OneOf) > 0 {
			for _, oneOf := range field.OneOf {
				fd := msg.FindFieldDescriptorByName(oneOf.Name)
				if fd == nil {
					return "", fmt.Errorf("can't find field descriptor by name '%s'", oneOf.Name)
				}
				oneOfDesc := fd.GetOneOf()
				if oneOfDesc == nil {
					return "", fmt.Errorf("file descriptor expected to be oneof '%s'", fd.GetFullyQualifiedName())
				}
				fdOneOf, value, err := msg.TryGetOneOfField(oneOfDesc)
				if err != nil {
					return "", fmt.Errorf("failed to find one of descriptor '%s': %w", oneOfDesc.GetFullyQualifiedName(), err)
				}
				if value == nil {
					break
				}
				tpl = strings.Repeat(" ", space) + field.Name + ": {%s: %s},\n"
				found := false
				for i := range field.OneOf {
					if field.OneOf[i].Name == fdOneOf.GetName() {
						field = field.OneOf[i]
						found = true
						break
					}
				}
				if !found {
					continue
				}
				protoValue = value
				break
			}
		} else {
			protoValue, err = msg.TryGetFieldByName(field.Name)
			if err != nil {
				return "", err
			}
		}
		if protoValue == nil {
			continue
		}
		if msgValue, ok := protoValue.(*dynamic.Message); ok && msgValue == nil {
			continue
		}
		v, err := f.asValue(field, protoValue, space)
		if err != nil {
			return "", err
		}
		if v == "" {
			continue
		}

		line := fmt.Sprintf(tpl, field.Name, v)
		buf.WriteString(line)
	}

	closeBracket := fmt.Sprintf("%s}", strings.Repeat(" ", space-2))
	buf.WriteString(closeBracket)
	return buf.String(), nil
}

func (f *Factory) asValue(field models.Field, value interface{}, space int) (string, error) {
	if field.Repeated {
		sliceValue, ok := value.([]interface{})
		if !ok {
			return "", fmt.Errorf("field %s expected as a collection, given '%s'", field.FullName, value)
		}
		sliceBuf := &strings.Builder{}
		sliceBuf.WriteString("[")
		for i, itemValue := range sliceValue {
			fieldCp := field
			fieldCp.Repeated = false
			v, err := f.asValue(fieldCp, itemValue, space+2)
			if err != nil {
				return "", err
			}
			sliceBuf.WriteString(v)
			if i != len(sliceValue)-1 {
				sliceBuf.WriteString(", ")
			}
		}
		sliceBuf.WriteString("]")
		return sliceBuf.String(), nil
	}

	var v string
	var err error
	switch field.Type {
	case models.DataTypeString:
		strV, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("field %s expected as string, given '%s'", field.FullName, value)
		}
		v = fmt.Sprintf(`'%s'`, strV)
	case models.DataTypeBool:
		boolV, ok := value.(bool)
		if !ok {
			return "", fmt.Errorf("field %s expected as boolean, given '%s'", field.FullName, value)
		}
		if boolV {
			v = "true"
		} else {
			v = "false"
		}
	case models.DataTypeInt32, models.DataTypeInt64, models.DataTypeUint32, models.DataTypeUint64:
		switch numV := value.(type) {
		case int32, int64, uint32, uint64:
			v = fmt.Sprintf(`%d`, numV)
		default:
			return "", fmt.Errorf("field %s expected as integer, given '%s'", field.FullName, value)
		}
	case models.DataTypeFloat32, models.DataTypeFloat64:
		switch numV := value.(type) {
		case float32, float64:
			v = fmt.Sprintf(`%g`, numV)
		default:
			return "", fmt.Errorf("field %s expected as float, given '%s'", field.FullName, value)
		}
	case models.DataTypeBytes:
		bytesV, ok := value.([]byte)
		if !ok {
			return "", fmt.Errorf("field %s expected as bytes, given '%s'", field.FullName, value)
		}
		v = fmt.Sprintf(`'%s'`, bytesV)
	case models.DataTypeEnum:
		if len(field.Enum) == 0 {
			return "", nil
		}
		numV, ok := value.(int32)
		if !ok {
			return "", fmt.Errorf("field %s expected as enum(int32), given '%s'", field.FullName, value)
		}
		v = fmt.Sprintf(`%d`, numV)
	case models.DataTypeDuration:
		durValue, ok := value.(*durationpb.Duration)
		if !ok {
			return "", fmt.Errorf("field %s expected as Duration, given '%s'", field.FullName, value)
		}
		v = fmt.Sprintf(`'%s'`, durValue.AsDuration().String())
	case models.DataTypeDate:
		timeValue, ok := value.(*timestamppb.Timestamp)
		if !ok {
			return "", fmt.Errorf("field %s expected as Timestamp, given '%s'", field.FullName, value)
		}
		v = fmt.Sprintf(`'%s'`, timeValue.AsTime().Format(time.RFC3339))
	case models.DataTypeStruct:
		if field.MapKey != nil && field.MapValue != nil {
			mapValue, ok := value.(map[any]any)
			if !ok {
				return "", fmt.Errorf("field %s expected as map, given '%s'", field.FullName, value)
			}
			mapBuf := &strings.Builder{}
			mapBuf.WriteString("{\n")
			for k, v := range mapValue {
				key, err := f.asValue(*field.MapKey, k, space+2)
				if err != nil {
					return "", nil
				}
				value, err := f.asValue(*field.MapValue, v, space+2)
				if err != nil {
					return "", nil
				}
				mapBuf.WriteString(fmt.Sprintf("%s%s: %s,\n", strings.Repeat(" ", space+2), key, value))
			}
			mapBuf.WriteString(strings.Repeat(" ", space))
			mapBuf.WriteString("}")
			v = mapBuf.String()
		} else {
			if value == nil || reflect.ValueOf(value).IsNil() {
				return "null", nil
			}

			link, ok := f.links[field.Message]
			if !ok {
				return "", fmt.Errorf("object %s not found", field.FullName)
			}

			m, ok := value.(*dynamic.Message)
			if !ok {
				protoMessage, ok := value.(protoiface.MessageV1)
				if !ok {
					return "", fmt.Errorf("unknown response type: %s", value)
				}
				m, err = dynamic.AsDynamicMessage(protoMessage)
				if err != nil {
					return "", fmt.Errorf("failed to cast to dynamic message: %w", err)
				}
			}

			v, err = f.asJs(link, m, space+2)
			if err != nil {
				return "", fmt.Errorf("field %s expected as map, given '%s'", field.FullName, value)
			}
		}
	default:
		return "", fmt.Errorf("data type is undefined: %s", field.Type)
	}

	return v, err
}
