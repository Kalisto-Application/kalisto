package definitions

import (
	"fmt"
	"kalisto/src/models"
	"sync"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

type Registry interface {
	Schema() (models.Spec, error)
	NewResponseMessage(methodFullName string) (*dynamic.Message, error)
	GetInputType(methodFullName string) (*desc.MessageDescriptor, error)
	GetOutputType(methodFullName string) (*desc.MessageDescriptor, error)
	MethodPath(methodFullName string) (string, error)
	Links() map[string]models.Message
}

type RegistryStore struct {
	mx    sync.RWMutex
	descs map[string]Registry
}

func NewRegistryStore() *RegistryStore {
	return &RegistryStore{descs: make(map[string]Registry)}
}

func (d *RegistryStore) Get(workspaceID string) (Registry, error) {
	d.mx.RLock()
	defer d.mx.RUnlock()

	r, ok := d.descs[workspaceID]
	if !ok {
		return nil, fmt.Errorf("failed to find descritors by workspace id==%s", workspaceID)
	}

	return r, nil
}

func (d *RegistryStore) Add(workspaceID string, r Registry) {
	d.mx.Lock()
	defer d.mx.Unlock()

	d.descs[workspaceID] = r
}
