package compiler

import (
	"errors"
	"fmt"
	"kalisto/src/definitions"
	"kalisto/src/models"
	"strings"
	"sync"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Registry struct {
	descriptors []*desc.FileDescriptor

	mx    sync.RWMutex
	links map[string]models.Message
}

func New(descriptors []*desc.FileDescriptor) *Registry {
	return &Registry{descriptors: descriptors, links: make(map[string]models.Message)}
}

var ErrFieldIsOneOf = errors.New("field is one of")

func (r *Registry) Schema() (spec models.Spec, err error) {
	if err := r.linkMessages(); err != nil {
		return spec, err
	}

	for _, fd := range r.descriptors {
		services := fd.GetServices()
		specServices := make([]models.Service, len(services))

		for i, service := range services {
			methods := service.GetMethods()
			specMethods := make([]models.Method, len(methods))

			for j, method := range methods {

				input := method.GetInputType().GetFullyQualifiedName()
				msg, ok := r.links[input]
				if !ok {
					return spec, fmt.Errorf("type %s not found: %w", input, err)
				}
				out := method.GetOutputType().GetFullyQualifiedName()
				responseMsg, ok := r.links[out]
				if !ok {
					return spec, fmt.Errorf("type %s not found: %w", out, err)
				}

				specMethods[j] = models.Method{
					Name:            method.GetName(),
					FullName:        method.GetFullyQualifiedName(),
					RequestMessage:  msg,
					ResponseMessage: responseMsg,
					Kind:            models.NewCommunicationKind(method.IsClientStreaming(), method.IsServerStreaming()),
					RequestExample:  definitions.MakeRequestExample(msg, r.Links()),
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

	spec.Links = r.Links()
	return spec, err
}

func (r *Registry) newField(fd *desc.FieldDescriptor, oneOfName string, set map[string]bool) (_ models.Field, err error) {
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
				key, err := r.newField(fd.GetMapKeyType(), "", set)
				if err != nil {
					return models.Field{}, err
				}
				mapKey = &key

				valueField, err := r.newField(fd.GetMapValueType(), "", set)
				if err != nil {
					return models.Field{}, err
				}
				mapValue = &valueField
			}

			dataType = models.DataTypeStruct
			message := fd.GetMessageType()
			linkKey := message.GetFullyQualifiedName()
			if _, ok := r.links[linkKey]; !ok {
				if err := r.linkMessageFields(message, linkKey); err != nil {
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
			oneOfField, err := r.newField(choice, oneOf.GetName(), set)
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

func (r *Registry) linkMessages() (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	protoregistry.GlobalTypes.RangeMessages(func(t protoreflect.MessageType) bool {
		mt, err := desc.WrapMessage(t.Descriptor())
		if err != nil {
			fmt.Println("failed to wrap message descriptor: ", string(t.Descriptor().FullName()), err.Error())
		}

		fullName := mt.GetFullyQualifiedName()
		if _, ok := r.links[fullName]; ok {
			return true
		}

		if err := r.linkMessageFields(mt, fullName); err != nil {
			fmt.Println("failed to link default msg: ", fullName, err.Error())
			return false
		}

		return true
	})

	for _, fd := range r.descriptors {
		mTypes := fd.GetMessageTypes()
		for _, mt := range mTypes {
			fullName := mt.GetFullyQualifiedName()
			if _, ok := r.links[fullName]; ok {
				continue
			}

			if err := r.linkMessageFields(mt, fullName); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Registry) linkMessageFields(mt *desc.MessageDescriptor, key string) error {
	mFields := mt.GetFields()
	fields := make([]models.Field, 0, len(mFields))
	r.links[key] = models.Message{
		Name:     mt.GetName(),
		FullName: mt.GetFullyQualifiedName(),
	}

	set := make(map[string]bool)
	for _, mf := range mFields {
		mField, err := r.newField(mf, "", set)
		if err != nil {
			if errors.Is(err, ErrFieldIsOneOf) {
				continue
			}
			return err
		}
		fields = append(fields, mField)
	}
	r.links[key] = models.Message{
		Name:     mt.GetName(),
		FullName: mt.GetFullyQualifiedName(),
		Fields:   fields,
	}

	return nil
}

func (r *Registry) FindMessage(name string) (*desc.MessageDescriptor, error) {
	for _, d := range r.descriptors {
		if m := d.FindMessage(name); m != nil {
			return m, nil
		}
	}

	return nil, fmt.Errorf("message not found")
}

func (r *Registry) FindMethod(full models.MethodName) (*desc.ServiceDescriptor, *desc.MethodDescriptor, error) {
	service, method := full.ServiceAndShort()
	for _, d := range r.descriptors {
		if s := d.FindService(service); s != nil {
			if m := s.FindMethodByName(method); m != nil {
				return s, m, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("method not found")
}

func (r *Registry) NewResponseMessage(methodFullName string) (*dynamic.Message, error) {
	_, md, err := r.FindMethod(models.MethodName(methodFullName))
	if err != nil {
		return nil, err
	}

	return dynamic.NewMessage(md.GetOutputType()), nil
}

// func (r *Registry) NewMessage(msg interface{}, name string) (interface{}, error) {
// 	msgType := msg.(*dynamic.Message).GetMessageDescriptor().FindFieldByName(name).GetMessageType()
// 	return dynamic.NewMessage(msgType), nil
// }

func (r *Registry) GetInputType(methodFullName string) (*desc.MessageDescriptor, error) {
	_, md, err := r.FindMethod(models.MethodName(methodFullName))
	if err != nil {
		return nil, err
	}

	return md.GetInputType(), nil
}

func (r *Registry) GetOutputType(methodFullName string) (*desc.MessageDescriptor, error) {
	_, md, err := r.FindMethod(models.MethodName(methodFullName))
	if err != nil {
		return nil, err
	}

	return md.GetOutputType(), nil
}

func (r *Registry) MethodPath(methodFullName string) (string, error) {
	sd, md, err := r.FindMethod(models.MethodName(methodFullName))
	if err != nil {
		return "", err
	}

	return "/" + sd.GetFullyQualifiedName() + "/" + md.GetName(), nil
}

// func (r *Registry) GetEnumByNumber(msg interface{}, name string, num int32) (int32, error) {
// 	v := msg.(*dynamic.Message).GetMessageDescriptor().FindFieldByName(name).GetMessageType().FindFieldByName(name).GetEnumType().FindValueByNumber(num)

// 	if v == nil {
// 		return 0, fmt.Errorf("enum not found")
// 	}
// 	return v.GetNumber(), nil
// }

// func (r *Registry) GetEnumByName(msg interface{}, name string, enum string) (int32, error) {
// 	v := msg.(*dynamic.Message).GetMessageDescriptor().FindFieldByName(name).GetMessageType().FindFieldByName(name).GetEnumType().FindValueByName(enum)

// 	if v == nil {
// 		return 0, fmt.Errorf("enum not found")
// 	}
// 	return v.GetNumber(), nil
// }

func (r *Registry) Links() map[string]models.Message {
	r.mx.RLock()
	defer r.mx.RUnlock()

	links := make(map[string]models.Message, len(r.links))
	for k, v := range r.links {
		links[k] = v
	}

	return links
}
