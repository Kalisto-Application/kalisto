package compiler

import (
	"fmt"
	"kalisto/src/models"
	"sync"

	"github.com/jhump/protoreflect/desc"
)

type Registry struct {
	Descriptors []*desc.FileDescriptor
}

func (r *Registry) FindMessage(name string) (*desc.MessageDescriptor, error) {
	for _, d := range r.Descriptors {
		if m := d.FindMessage(name); m != nil {
			return m, nil
		}
	}

	return nil, fmt.Errorf("message not found")
}

func (r *Registry) FindMethod(full models.MethodName) (*desc.ServiceDescriptor, *desc.MethodDescriptor, error) {
	service, method := full.ServiceAndShort()
	for _, d := range r.Descriptors {
		if s := d.FindService(service); s != nil {
			if m := s.FindMethodByName(method); m != nil {
				return s, m, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("method not found")
}

type Descritors struct {
	mx    sync.RWMutex
	descs map[string]*Registry
}

func NewProtoRegistry() *Descritors {
	return &Descritors{descs: make(map[string]*Registry)}
}

func (d *Descritors) Get(workspaceID string) (*Registry, error) {
	d.mx.RLock()
	defer d.mx.RUnlock()

	r, ok := d.descs[workspaceID]
	if !ok {
		return nil, fmt.Errorf("failed to find descritors by workspace id==%s", workspaceID)
	}

	return r, nil
}

func (d *Descritors) Add(workspaceID string, r *Registry) {
	d.mx.Lock()
	defer d.mx.Unlock()

	d.descs[workspaceID] = r
}
