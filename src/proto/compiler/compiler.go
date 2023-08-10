package compiler

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

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

func (c *FileCompiler) Compile(path string, filenames []string, bufDirs []string) (*Registry, error) {
	if len(filenames) == 0 {
		return nil, ErrNoFiles
	}

	paths, filenames, err := c.cutWorkspaceDirs(path, filenames, bufDirs)
	if err != nil {
		return nil, err
	}

	parser := protoparse.Parser{ImportPaths: paths}
	descriptors, err := parser.ParseFiles(filenames...)
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

func (c *FileCompiler) cutWorkspaceDirs(path string, filenames []string, bufDirs []string) ([]string, []string, error) {
	filenames, err := protoparse.ResolveFilenames([]string{path}, filenames...)
	if err != nil {
		return nil, nil, fmt.Errorf("compiler: failed to resolve files: %w", err)
	}

	if len(bufDirs) == 0 {
		return []string{path}, filenames, nil
	}
	pathSet := make(map[string]struct{})

	for i, name := range filenames {
		for _, bufDir := range bufDirs {
			newName, ok := strings.CutPrefix(name, bufDir)
			if ok {
				pathSet[filepath.Join(path, bufDir)] = struct{}{}
				// it cuts directory separator (like "/") OS independently
				newName = newName[1:]
				filenames[i] = newName
				break
			}
		}
	}

	paths := make([]string, 0, len(bufDirs))
	for bufPath := range pathSet {
		paths = append(paths, bufPath)
	}

	return paths, filenames, nil
}
