package compiler

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
)

type Compiler struct {
}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(filenames []string) (*Registry, error) {
	docs := make([]*openapi3.T, 0)

	for _, filename := range filenames {
		doc, err := (&openapi3.Loader{Context: context.Background(), IsExternalRefsAllowed: true}).LoadFromFile(filename)
		if err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}

	return NewRegistry(docs), nil
}
