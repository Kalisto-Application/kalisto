package interpreter

import (
	"fmt"
	"kalisto/src/models"
	"math"
	"strconv"

	"github.com/dop251/goja"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

func CreateMessageFromScript(script string, desc *desc.MessageDescriptor, spec models.Spec, serviceName, methodName string) (*dynamic.Message, error) {
	defer func() {
		x := recover()
		fmt.Println(x)
	}()

	message, err := spec.FindInputMessage(serviceName, methodName)
	if err != nil {
		return nil, err
	}

	script = fmt.Sprintf(`(() => {
		return %s
	  })()`, script)

	vm := goja.New()
	val, err := vm.RunString(script)
	if err != nil {
		return nil, fmt.Errorf("interpretator: failed to run script: %w", err)
	}
	if val == nil {
		return nil, nil
	}

	m, ok := val.Export().(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("interpretator: failed to convert script to map: %w", err)
	}

	return newMessage(desc, spec, m, message)
}

func newMessage(desc *desc.MessageDescriptor, spec models.Spec, m map[string]interface{}, message models.Message) (*dynamic.Message, error) {
	resultMessage := dynamic.NewMessage(desc)

	for k, v := range m {
		field, err := message.FindField(k)
		if err != nil {
			return nil, err
		}

		v, err = castValue(desc, spec, field, v)
		if err != nil {
			return nil, err
		}
		if err := resultMessage.TrySetFieldByName(k, v); err != nil {
			return nil, err
		}
	}

	return resultMessage, nil
}

func castValue(desc *desc.MessageDescriptor, spec models.Spec, f models.Field, v interface{}) (interface{}, error) {
	if v == nil {
		return v, nil
	}

	if f.MapKey != nil {
		val, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to cast map value, expected object")
		}
		mDesc := desc.FindFieldByName(f.Name).GetMessageType()
		ret := make([]*dynamic.Message, 0, len(val))
		for k, v := range val {
			msg := dynamic.NewMessage(mDesc)
			key, err := castValue(mDesc, spec, *f.MapKey, k)
			if err != nil {
				return nil, err
			}

			value, err := castValue(mDesc, spec, *f.MapValue, v)
			if err != nil {
				return nil, err
			}

			if err := msg.TrySetFieldByName("key", key); err != nil {
				return nil, err
			}
			if err := msg.TrySetFieldByName("value", value); err != nil {
				return nil, err
			}
			ret = append(ret, msg)
		}

		return ret, nil
	}

	switch f.Type {
	case models.DataTypeBool:
		if strV, ok := v.(string); ok {
			boolV, err := strconv.ParseBool(strV)
			if err != nil {
				return nil, fmt.Errorf("expected bool: %w", err)
			}
			return boolV, nil
		}
		return v, nil
	case models.DataTypeInt32:
		var val int32
		if intV, ok := v.(int64); ok {
			val = int32(intV)
		}
		if floatV, ok := v.(float64); ok {
			val = int32(floatV)
		}
		if strV, ok := v.(string); ok {
			intV, err := strconv.ParseInt(strV, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("expected int32: %w", err)
			}
			val = int32(intV)
		}
		if val < math.MinInt32 || val > math.MaxInt32 {
			return nil, fmt.Errorf("value is out of range")
		}
		return val, nil
	case models.DataTypeInt64:
		var val int64
		if intV, ok := v.(int64); ok {
			val = int64(intV)
		}
		if floatV, ok := v.(float64); ok {
			val = int64(floatV)
		}
		if val < math.MinInt64 || val > math.MaxInt64 {
			return nil, fmt.Errorf("value is out of range")
		}
		if strV, ok := v.(string); ok {
			intV, err := strconv.ParseInt(strV, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("expected int64: %w", err)
			}
			val = intV
		}
		return val, nil
	case models.DataTypeUint32:
		var val uint32
		if intV, ok := v.(int64); ok {
			if intV < 0 || intV > math.MaxUint32 {
				return nil, fmt.Errorf("value is out of range")
			}
			val = uint32(intV)
		}
		if floatV, ok := v.(float64); ok {
			if floatV < 0 || floatV > math.MaxUint32 {
				return nil, fmt.Errorf("value is out of range")
			}
			val = uint32(floatV)
		}
		if strV, ok := v.(string); ok {
			intV, err := strconv.ParseUint(strV, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("expected int64: %w", err)
			}
			val = uint32(intV)
		}
		return val, nil
	case models.DataTypeUint64:
		var val uint64
		if intV, ok := v.(int64); ok {
			if intV < 0 {
				return nil, fmt.Errorf("value is out of range")
			}
			val = uint64(intV)
		}
		if floatV, ok := v.(float64); ok {
			if floatV < 0 || floatV > math.MaxUint64 {
				return nil, fmt.Errorf("value is out of range")
			}
			val = uint64(floatV)
		}
		if strV, ok := v.(string); ok {
			intV, err := strconv.ParseUint(strV, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("expected int64: %w", err)
			}
			val = intV
		}
		return val, nil
	case models.DataTypeFloat32:
		if intV, ok := v.(int64); ok {
			return float32(intV), nil
		}
		if floatV, ok := v.(float64); ok {
			if floatV > math.MaxFloat32 || floatV < math.SmallestNonzeroFloat32 {
				return nil, fmt.Errorf("value is out of range")
			}
			return float32(floatV), nil
		}
	case models.DataTypeFloat64:
		if intV, ok := v.(int64); ok {
			return float64(intV), nil
		}
		if floatV, ok := v.(float64); ok {
			return float64(floatV), nil
		}
	case models.DataTypeString:
		return v, nil
	case models.DataTypeBytes:
		strV, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("value must be string")
		}
		return []byte(strV), nil
	case models.DataTypeEnum:
		intV, ok := v.(int64)
		if !ok {
			return nil, fmt.Errorf("enum value is out of range")
		}
		return int32(intV), nil
	case models.DataTypeStruct:
		val, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("can not cast message type, given %T", v)
		}

		fieldDesc := desc.FindFieldByName(f.Name)
		if fieldDesc == nil {
			return nil, fmt.Errorf("descriptor field=%s not found", f.Name)
		}
		link, ok := spec.Links[f.Message]
		if !ok {
			return nil, fmt.Errorf("failed to find message link=%s", f.Message)
		}
		return newMessage(fieldDesc.GetMessageType(), spec, val, link)
	// case models.DataTypeMap:
	// 	if f.MapKey == nil {
	// 		return nil, fmt.Errorf("map must have collection key, found nil")
	// 	}

	// 	// keyMessage := newField(f.CollectionKey)
	// 	// valueMessage := newField(f.)

	// return v, nil
	// case models.DataTypeArray:
	// return v, nil
	case models.DataTypeDate:
		return v, nil
	default:
		return nil, fmt.Errorf("undefined type=%s", f.Type)
	}
	return v, nil
}

func messageFromScript() (*dynamic.Message, error) {
	return nil, nil
}
