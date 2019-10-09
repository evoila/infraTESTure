package parser

import (
	"github.com/evoila/infraTESTure/config"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func GetMethodNames(tests []config.Test, path string) ([]string, error) {
	var methodNames []string

	fileSet := token.NewFileSet() // positions are relative to fileSet
	file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, test := range tests {
		for _, decl := range file.Decls {
			if fun, ok := decl.(*ast.FuncDecl); ok {
				if fun.Doc != nil && contains("@"+test.Name, fun.Doc.List) {
					methodNames = append(methodNames, fun.Name.Name)
				}
			}
		}
	}

	return methodNames, nil
}

func contains(s string, e []*ast.Comment) bool {
	for _, b := range e {
		if strings.Contains(strings.ToLower(b.Text), strings.ToLower(s)) {
			return true
		}
	}
	return false
}