package render

import (
	"fmt"
	"kalisto/src/models"
	"reflect"
	"strings"
	"time"

	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Renderer struct {
	links map[string]models.Message
}

func New(links map[string]models.Message) *Renderer {
	return &Renderer{links: links}
}

func (f *Renderer) MessageAsJsString(spec models.Message, msg *dynamic.Message) (string, error) {
	msgMap, err := f.MessageAsMap(spec, msg)
	if err != nil {
		return "", err
	}

	return f.MapAsJs(spec, msgMap)
}

func (f *Renderer) MessageAsMap(spec models.Message, msg *dynamic.Message) (map[string]interface{}, error) {
	res := make(map[string]interface{}, len(spec.Fields))

	for _, field := range spec.Fields {

		protoValue, mutatedField, err := f.getProtoValue(field, msg)
		if err != nil {
			return res, err
		}
		if protoValue == nil {
			continue
		}

		val, err := f.asMapValue(mutatedField, protoValue)
		if err != nil {
			return nil, err
		}

		if field.Type == models.DataTypeOneOf {
			val = map[string]interface{}{
				mutatedField.Name: val,
			}
		}
		res[field.Name] = val
	}

	return res, nil
}

func (f *Renderer) MapAsJs(m models.Message, val map[string]interface{}) (string, error) {
	return f.mapAsJs(m, val, 2)
}

func (f *Renderer) mapAsJs(m models.Message, val map[string]interface{}, space int) (string, error) {
	var buf strings.Builder
	buf.WriteString("{\n")

	for _, field := range m.Fields {
		value := val[field.Name]
		tpl := strings.Repeat(" ", space) + "%s: %s,\n"
		v, err := f.mapAsJsValue(value, space)
		if err != nil {
			return "", fmt.Errorf("failed to present a map as js: %w", err)
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

func (f *Renderer) asMapValue(field models.Field, value interface{}) (interface{}, error) {
	if field.Repeated {
		sliceValue, ok := value.([]interface{})
		if !ok {
			return "", fmt.Errorf("field %s expected as a collection, given '%s'", field.FullName, value)
		}
		values := make([]interface{}, len(sliceValue))
		for i, itemValue := range sliceValue {
			fieldCp := field
			fieldCp.Repeated = false
			v, err := f.asMapValue(fieldCp, itemValue)
			if err != nil {
				return "", err
			}
			values[i] = v
		}
		return values, nil
	}

	switch field.Type {
	case models.DataTypeString:
		strV, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("field %s expected as string, given '%s'", field.FullName, value)
		}
		return strV, nil
	case models.DataTypeBool:
		boolV, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("field %s expected as boolean, given '%s'", field.FullName, value)
		}
		return boolV, nil
	case models.DataTypeInt32, models.DataTypeInt64, models.DataTypeUint32, models.DataTypeUint64:
		switch numV := value.(type) {
		case int32, int64, uint32, uint64:
			return numV, nil
		default:
			return nil, fmt.Errorf("field %s expected as integer, given '%s'", field.FullName, value)
		}
	case models.DataTypeFloat32, models.DataTypeFloat64:
		switch numV := value.(type) {
		case float32, float64:
			return numV, nil
		default:
			return nil, fmt.Errorf("field %s expected as float, given '%s'", field.FullName, value)
		}
	case models.DataTypeBytes:
		bytesV, ok := value.([]byte)
		if !ok {
			return nil, fmt.Errorf("field %s expected as bytes, given '%s'", field.FullName, value)
		}
		return bytesV, nil
	case models.DataTypeEnum:
		if len(field.Enum) == 0 {
			return 0, nil
		}
		numV, ok := value.(int32)
		if !ok {
			return nil, fmt.Errorf("field %s expected as enum(int32), given '%s'", field.FullName, value)
		}
		return numV, nil
	case models.DataTypeDuration:
		durValue, ok := value.(*durationpb.Duration)
		if !ok {
			return nil, fmt.Errorf("field %s expected as Duration, given '%s'", field.FullName, value)
		}
		return int64(durValue.AsDuration()), nil
	case models.DataTypeDate:
		timeValue, ok := value.(*timestamppb.Timestamp)
		if !ok {
			return nil, fmt.Errorf("field %s expected as Timestamp, given '%s'", field.FullName, value)
		}
		return timeValue.AsTime(), nil
	case models.DataTypeStruct:
		if field.MapKey != nil && field.MapValue != nil {
			values := make(map[interface{}]interface{})
			mapValue, ok := value.(map[any]any)
			if !ok {
				return nil, fmt.Errorf("field %s expected as map, given '%s'", field.FullName, value)
			}
			for k, v := range mapValue {
				key, err := f.asMapValue(*field.MapKey, k)
				if err != nil {
					return nil, err
				}
				value, err := f.asMapValue(*field.MapValue, v)
				if err != nil {
					return nil, err
				}
				values[key] = value
			}
			return values, nil
		} else {
			if value == nil || reflect.ValueOf(value).IsNil() {
				return nil, nil
			}

			link, ok := f.links[field.Message]
			if !ok {
				return "", fmt.Errorf("object %s not found", field.FullName)
			}
			var err error
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

			v, err := f.MessageAsMap(link, m)
			if err != nil {
				return "", fmt.Errorf("field %s expected as map, given '%s'", field.FullName, value)
			}
			return v, nil
		}
	default:
		return nil, fmt.Errorf("data type is undefined: %s", field.Type)
	}
}

func (f *Renderer) mapAsJsValue(value interface{}, space int) (string, error) {
	switch val := value.(type) {
	case []interface{}:
		sliceBuf := &strings.Builder{}
		sliceBuf.WriteString("[")

		for i, itemValue := range val {
			v, err := f.mapAsJsValue(itemValue, space+2)
			if err != nil {
				return "", err
			}
			if v == "" {
				continue
			}
			sliceBuf.WriteString(v)
			if i != len(val)-1 {
				sliceBuf.WriteString(", ")
			}
		}
		sliceBuf.WriteString("]")
		return sliceBuf.String(), nil

	case string:
		return fmt.Sprintf(`'%s'`, val), nil
	case []byte:
		return fmt.Sprintf(`'%s'`, string(val)), nil
	case bool:
		v := "false"
		if val {
			v = "true"
		}
		return v, nil
	case int32, int64, uint32, uint64:
		return fmt.Sprintf(`%d`, val), nil
	case float32, float64:
		return fmt.Sprintf(`%g`, val), nil
	case time.Duration:
		return fmt.Sprintf(`'%s'`, val.String()), nil
	case time.Time:
		return fmt.Sprintf(`'%s'`, val.Format(time.RFC3339)), nil
	case map[interface{}]interface{}:
		mapBuf := &strings.Builder{}
		mapBuf.WriteString("{\n")
		for k, v := range val {
			key, err := f.mapAsJsValue(k, space+2)
			if err != nil {
				return "", nil
			}
			if key == "" {
				continue
			}
			value, err := f.mapAsJsValue(v, space+2)
			if err != nil {
				return "", nil
			}
			if value == "" {
				continue
			}
			mapBuf.WriteString(fmt.Sprintf("%s%s: %s,\n", strings.Repeat(" ", space+2), key, value))
		}
		mapBuf.WriteString(strings.Repeat(" ", space))
		mapBuf.WriteString("}")
		return mapBuf.String(), nil
	case map[string]interface{}:
		mapBuf := &strings.Builder{}
		mapBuf.WriteString("{\n")
		for k, v := range val {
			value, err := f.mapAsJsValue(v, space+2)
			if err != nil {
				return "", nil
			}
			if value == "" {
				continue
			}
			mapBuf.WriteString(fmt.Sprintf("%s%s: %s,\n", strings.Repeat(" ", space+2), k, value))
		}
		mapBuf.WriteString(strings.Repeat(" ", space))
		mapBuf.WriteString("}")
		return mapBuf.String(), nil
	case nil:
		return "", nil
	}

	return "", fmt.Errorf("unknown type: %v", value)
}

func (f *Renderer) getProtoValue(field models.Field, msg *dynamic.Message) (interface{}, models.Field, error) {
	var protoValue interface{}
	var err error

	if len(field.OneOf) > 0 {
		for _, oneOf := range field.OneOf {
			fd := msg.FindFieldDescriptorByName(oneOf.Name)
			if fd == nil {
				return "", field, fmt.Errorf("can't find field descriptor by name '%s'", oneOf.Name)
			}
			oneOfDesc := fd.GetOneOf()
			if oneOfDesc == nil {
				return "", field, fmt.Errorf("file descriptor expected to be oneof '%s'", fd.GetFullyQualifiedName())
			}
			fdOneOf, value, err := msg.TryGetOneOfField(oneOfDesc)
			if err != nil {
				return "", field, fmt.Errorf("failed to find one of descriptor '%s': %w", oneOfDesc.GetFullyQualifiedName(), err)
			}
			if value == nil {
				break
			}
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
			return "", field, err
		}
	}

	return protoValue, field, err
}
