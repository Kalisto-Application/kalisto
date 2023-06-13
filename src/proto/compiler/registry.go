package compiler

import (
	"github.com/jhump/protoreflect/desc"
)

type Registry struct {
	Descriptors []*desc.FileDescriptor
}
