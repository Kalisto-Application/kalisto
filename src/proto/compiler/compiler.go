package compiler

import (
	"errors"
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"

	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	ErrNoFiles = errors.New("compiler: files are required to compile proto")
)

type FileCompiler struct {
	seen map[string]struct{}
}

func NewFileCompiler() *FileCompiler {
	return &FileCompiler{seen: make(map[string]struct{})}
}

func (c *FileCompiler) Compile(paths, filenames []string) (*Registry, error) {
	if len(filenames) == 0 {
		return nil, ErrNoFiles
	}

	files, err := protoparse.ResolveFilenames(paths, filenames...)
	if err != nil {
		return nil, fmt.Errorf("compiler: failed to resolve files: %w", err)
	}

	parser := protoparse.Parser{ImportPaths: paths}
	descriptors, err := parser.ParseFiles(files...)
	if err != nil {
		return nil, fmt.Errorf("compiler: failed to parse files: %w", err)
	}

	return &Registry{Descriptors: descriptors}, nil
	// fdset := &descriptorpb.FileDescriptorSet{}
	// for _, fd := range descriptors {
	// 	fdset.File = append(fdset.File, c.walk(fd)...)
	// }

	// return protodesc.NewFiles(fdset)
}

func (c *FileCompiler) walk(fd *desc.FileDescriptor) []*descriptorpb.FileDescriptorProto {
	descriptorsProto := []*descriptorpb.FileDescriptorProto{}

	key := fd.GetName() + fd.GetPackage()
	if _, ok := c.seen[key]; ok {
		return descriptorsProto
	}
	c.seen[key] = struct{}{}
	descriptorsProto = append(descriptorsProto, fd.AsFileDescriptorProto())

	for _, dep := range fd.GetDependencies() {
		deps := c.walk(dep)
		descriptorsProto = append(descriptorsProto, deps...)
	}

	return descriptorsProto
}
