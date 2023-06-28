package interpreter

import (
	"fmt"
	"kalisto/src/models"

	"github.com/dop251/goja"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

func CreateMessageFromScript(script string, desc *desc.MessageDescriptor, spec models.Spec, serviceName, methodName string) (*dynamic.Message, error) {
	script = fmt.Sprintf(`(() => {
		return %s
	  })()`, script)

	vm := goja.New()
	val, err := vm.RunString(script)
	if err != nil {
		return nil, fmt.Errorf("interpretator: failed to run script: %w", err)
	}

	m, ok := val.Export().(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("interpretator: failed to convert script to map: %w", err)
	}

	resultMessage := dynamic.NewMessage(desc)
	for k, v := range m {
		// messageSpec, err := spec.FindInputMessage()
		resultMessage.SetFieldByName(k, v)
	}

	return resultMessage, nil
}

func messageFromScript() (*dynamic.Message, error) {
	return nil, nil
}
