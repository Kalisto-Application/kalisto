package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"google.golang.org/protobuf/types/descriptorpb"
	"gopkg.in/yaml.v3"
)

type Registry struct {
	Descriptors []*desc.FileDescriptor
}

type bufWork struct {
	Directories []string `yaml:"directories"`
}

type ProtoSearchResult struct {
	AbsoluteDirPath    string
	RelativeProtoPaths []string
	BufDirs            []string
}

// SearchProtoFiles function will find all .proto files by the given path.
func SearchProtoFiles(path string) (ProtoSearchResult, error) {
	result := ProtoSearchResult{}

	// Check if the path is absolute
	if !filepath.IsAbs(path) {
		return result, errors.New("path must be absolute")
	}
	result.BufDirs = readBufWorkDirs(path)

	// Check if the path is a directory or a file
	info, err := os.Stat(path)
	if err != nil {
		return result, err
	}

	// This is a file
	if !info.IsDir() {
		if !strings.HasSuffix(info.Name(), ".proto") {
			return result, errors.New("chosen file is not a proto file")
		}

		result.AbsoluteDirPath = filepath.Dir(path)
		result.RelativeProtoPaths = []string{filepath.Base(path)}

		return result, nil
	}

	// This is a directory, find all .proto files recursively
	result.AbsoluteDirPath = path
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".proto") {
			rel, err := filepath.Rel(result.AbsoluteDirPath, path)
			if err != nil {
				return err
			}
			result.RelativeProtoPaths = append(result.RelativeProtoPaths, rel)
		}
		return nil
	})

	if err != nil {
		return result, err
	}

	if len(result.RelativeProtoPaths) == 0 {
		return result, fmt.Errorf("no proto files")
	}

	return result, nil
}

func readBufWorkDirs(path string) []string {
	f, err := os.Open(filepath.Join(path, "buf.work.yaml"))
	if err != nil {
		f, err = os.Open(filepath.Join(path, "buf.work.yml"))
		if err != nil {
			return nil
		}
	}

	var buf bufWork
	if err := yaml.NewDecoder(f).Decode(&buf); err != nil {
		return nil
	}

	return buf.Directories
}

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

func main() {
	path := os.Args[1]
	protoFiels, err := SearchProtoFiles(path)
	if err != nil {
		log.Fatalln("err on search proto", err.Error())
	}
	fmt.Println("proto files: ", protoFiels)
	comp := NewFileCompiler()
	reg, err := comp.Compile(protoFiels.AbsoluteDirPath, protoFiels.RelativeProtoPaths, protoFiels.BufDirs)
	if err != nil {
		log.Fatalln("err on compile files: ", err)
	}

	for _, desc := range reg.Descriptors {
		for _, service := range desc.GetServices() {
			fmt.Println("service: ", service.GetFullyQualifiedName())
			for _, method := range service.GetMethods() {
				fmt.Println("method: ", method.GetFullyQualifiedName())
			}
		}
	}
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
