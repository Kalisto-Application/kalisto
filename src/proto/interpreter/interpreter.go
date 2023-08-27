package interpreter

import (
	"context"
	"errors"
	"fmt"
	"kalisto/src/models"
	"kalisto/src/proto/client"
	"kalisto/src/proto/compiler"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type Interpreter struct {
	vars string
}

func NewInterpreter(vars string) *Interpreter {
	return &Interpreter{vars: vars}
}

func (ip *Interpreter) CreateMessageFromScript(script string, desc *desc.MessageDescriptor, spec models.Spec, message models.Message) (*dynamic.Message, error) {
	m, err := ip.exportValue(script, "body.js")
	if err != nil {
		return nil, err
	}

	return newMessage(desc, spec, m, message)
}

func (ip *Interpreter) CreateMetadata(script string) (metadata.MD, error) {
	m, err := ip.exportValue(script, "header.js")
	if err != nil {
		return nil, err
	}

	return newMeta(m)
}

func (ip *Interpreter) Raw(script string) (map[string]interface{}, error) {
	return ip.exportValue(script, "script.js")
}

type requestFunc func(obj interface{}) (interface{}, error)
type apiError string

func (a apiError) Error() string {
	return string(a)
}

func (ip *Interpreter) RunScript(ctx context.Context, script string, spec models.Spec, reg *compiler.Registry, client *client.Client) (*dynamic.Message, error) {
	vm := goja.New()
	if ip.vars != "" {
		globalScript := fmt.Sprintf("g = %s;", ip.vars)
		if _, err := vm.RunScript("global.js", globalScript); err != nil {
			return nil, ip.mapErr(vm, err)
		}
	}

	for _, service := range spec.Services {
		jsService := vm.NewObject()
		for _, method := range service.Methods {
			sd, md, err := reg.FindMethod(models.MethodName(method.FullName))
			if err != nil {
				return nil, err
			}
			jsService.Set(method.Name, newJsFunc(ctx, vm, sd, md, spec, method, client))
		}

		if err := vm.Set(service.Name, jsService); err != nil {
			return nil, err
		}
	}

	value, err := vm.RunString(script)
	if err != nil {
		return nil, ip.mapErr(vm, err)
	}
	exported := value.Export()
	response, ok := exported.(*dynamic.Message)
	if !ok {
		return nil, fmt.Errorf("expected grpc message as a result")
	}
	return response, nil
}

func newJsFunc(ctx context.Context, vm *goja.Runtime, sd *desc.ServiceDescriptor, md *desc.MethodDescriptor, spec models.Spec, method models.Method, client *client.Client) requestFunc {
	return func(obj interface{}) (interface{}, error) {
		ctx = context.Background()
		var msg *dynamic.Message
		var err error
		switch obj := obj.(type) {
		case map[string]interface{}:
			msg, err = newMessage(md.GetInputType(), spec, obj, method.RequestMessage)
			if err != nil {
				return nil, err
			}
		case *dynamic.Message:
			msg = obj
		default:
			return nil, models.JsTypeError
		}

		resp := dynamic.NewMessage(md.GetOutputType())
		responseMeta := metadata.MD{}
		apiErr, err := client.Invoke(ctx, "/"+sd.GetFullyQualifiedName()+"/"+md.GetName(), msg, resp, &responseMeta)
		if err != nil {
			return nil, err
		}
		if apiErr != "" {
			return nil, apiError(apiErr)
		}

		return resp, nil
	}
}

func (ip *Interpreter) exportValue(script, name string) (map[string]interface{}, error) {
	for {
		if !strings.HasPrefix(script, "\n") {
			break
		}

		script = strings.TrimPrefix(script, "\n")
	}
	for {
		if !strings.HasSuffix(script, "\n") {
			break
		}

		script = strings.TrimSuffix(script, "\n")
	}

	script = fmt.Sprintf(`(()=> {
		return %s;
	})()`, script)

	vm := goja.New()
	if ip.vars != "" {
		globalScript := fmt.Sprintf("g = %s;", ip.vars)
		if _, err := vm.RunScript("global.js", globalScript); err != nil {
			return nil, ip.mapErr(vm, err)
		}
	}
	val, err := vm.RunScript(name, script)
	if err != nil {
		return nil, ip.mapErr(vm, err)
	}
	if val == nil || val.ExportType() == nil {
		return nil, nil
	}
	m, ok := val.Export().(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("interpretator: failed to convert script to map: %w", err)
	}

	return m, nil
}

func (ip *Interpreter) mapErr(vm *goja.Runtime, err error) error {
	var exc *goja.Exception
	if errors.As(err, &exc) {
		obj, ok := exc.Value().Export().(map[string]interface{})
		if ok && obj["value"] != nil && errors.Is(obj["value"].(error), models.JsTypeError) {
			return models.ErrorSyntax(models.JsTypeError.Error())
		}
		return models.ErrorSyntax(exc.Error())
	}
	return fmt.Errorf("interpretator: failed to run script: %w", err)
}

func newMessage(desc *desc.MessageDescriptor, spec models.Spec, m map[string]interface{}, message models.Message) (*dynamic.Message, error) {
	resultMessage := dynamic.NewMessage(desc)

	for k, v := range m {
		field, err := message.FindField(k)
		if err != nil {
			// TODO: warn the field is unused
			continue
		}

		value, err := castValue(desc, spec, field, v)
		if err != nil {
			return nil, err
		}
		if field.Type == models.DataTypeOneOf {
			for key := range v.(map[string]interface{}) {
				k = key
			}
		}

		if err := resultMessage.TrySetFieldByName(k, value); err != nil {
			return nil, err
		}
	}

	return resultMessage, nil
}

func newMeta(vals map[string]interface{}) (metadata.MD, error) {
	metaMap := make(map[string][]string, len(vals))
	for k, v := range vals {
		switch val := v.(type) {
		case string:
			metaMap[k] = []string{val}
		case []string:
			metaMap[k] = val
		case []interface{}:
			valStr := make([]string, len(val))
			for i := range val {
				vStr, ok := val[i].(string)
				if !ok {
					return nil, fmt.Errorf("header value must be string or string[]")
				}
				valStr[i] = vStr
			}
			metaMap[k] = valStr
		default:
			return nil, fmt.Errorf("only string and []string values are allowed in headers")
		}
	}

	return metadata.MD(metaMap), nil
}

func castValue(desc *desc.MessageDescriptor, spec models.Spec, f models.Field, v interface{}) (interface{}, error) {
	if v == nil {
		return v, nil
	}

	if f.Repeated {
		val, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to cast repeated value, expected array")
		}
		fCopy := f
		fCopy.Repeated = false
		ret := make([]interface{}, 0, len(val))
		for _, it := range val {
			casted, err := castValue(desc, spec, fCopy, it)
			if err != nil {
				return nil, err
			}
			ret = append(ret, casted)
		}
		return ret, nil
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
		switch v := v.(type) {
		case int64:
			fieldDesc := desc.FindFieldByName(f.Name)
			if fieldDesc == nil {
				return nil, fmt.Errorf("descriptor field=%s not found", f.Name)
			}
			enumValue := fieldDesc.GetEnumType().FindValueByNumber(int32(v))
			if enumValue == nil {
				return nil, models.ErrorSyntax(fmt.Sprintf("'%s': %d:  enum value is out of range", f.Name, v))
			}
			return enumValue.GetNumber(), nil
		case string:
			fieldDesc := desc.FindFieldByName(f.Name)
			if fieldDesc == nil {
				return nil, fmt.Errorf("descriptor field=%s not found", f.Name)
			}
			enumValue := fieldDesc.GetEnumType().FindValueByName(v)
			if enumValue == nil {
				return nil, models.ErrorSyntax(fmt.Sprintf("'%s': '%s':  enum value is out of range", f.Name, v))
			}
			return enumValue.GetNumber(), nil
		default:
			return nil, fmt.Errorf("enum value must be integer or string")
		}
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
	case models.DataTypeOneOf:
		val, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("can not cast message type, given %T", v)
		}
		if len(val) > 1 {
			return nil, fmt.Errorf("oneof must contain only one field")
		}
		for key, value := range val {
			oneOf, err := f.FindOneofByName(key)
			if err != nil {
				return nil, err
			}
			return castValue(desc, spec, oneOf, value)
		}

	case models.DataTypeDate:
		m, err := makeKnownMessage("google.protobuf.Timestamp")
		if err != nil {
			return nil, err
		}
		switch v := v.(type) {
		case int64:
			t := time.UnixMilli(v).UTC()
			if err := timeMessageFromNanos(t.UnixNano(), m); err != nil {
				return nil, err
			}
			return m, nil
		case string:
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, err
			}
			if err := timeMessageFromNanos(t.UnixNano(), m); err != nil {
				return nil, err
			}
			return m, nil
		case time.Time:
			if err := timeMessageFromNanos(v.UnixNano(), m); err != nil {
				return nil, err
			}
			return m, nil
		default:
			return nil, fmt.Errorf("date must be nanoseconds (number), timestamp (string) or Date instance")
		}
	case models.DataTypeDuration:
		m, err := makeKnownMessage("google.protobuf.Duration")
		if err != nil {
			return nil, err
		}
		switch v := v.(type) {
		case float64:
			if v > math.MaxInt64 {
				return nil, fmt.Errorf("max value of duration is int64: %d", math.MaxInt64)
			}
			if err := timeMessageFromNanos(int64(v), m); err != nil {
				return nil, err
			}
			return m, nil
		case int64:
			if err := timeMessageFromNanos(v, m); err != nil {
				return nil, err
			}
			return m, nil
		case string:
			d, err := ParseDuration(v)
			if err != nil {
				return nil, err
			}
			if err := timeMessageFromNanos(int64(d), m); err != nil {
				return nil, err
			}
			return m, nil
		default:
			return nil, fmt.Errorf("duration must be integer or string")
		}
	default:
		return nil, fmt.Errorf("undefined type=%s", f.Type)
	}
	return v, nil
}

func makeKnownMessage(name string) (*dynamic.Message, error) {
	mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(name))
	if err != nil {
		return nil, fmt.Errorf("failed to get descriptor: %w", err)
	}
	d, err := desc.WrapMessage(mt.Descriptor())
	if err != nil {
		return nil, fmt.Errorf("failed to wrap descriptor")
	}
	return dynamic.NewMessage(d), nil
}

func timeMessageFromNanos(nanos int64, m *dynamic.Message) error {
	secs := nanos / 1e9
	nanos -= secs * 1e9
	if err := m.TrySetFieldByName("seconds", secs); err != nil {
		return err
	}
	if err := m.TrySetFieldByName("nanos", int32(nanos)); err != nil {
		return err
	}
	return nil
}
