package parser

import (
	"github.com/evoila/infraTESTure/config"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
	"unicode/utf8"
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

func GetAnnotations(path string) ([]string, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	comments := ast.NewCommentMap(fileSet, file, file.Comments).Comments()

	for _, comment := range comments {
		if strings.HasPrefix(comment.Text(), "@") {
			log.Printf("│ \t├── %v", trim(comment.Text()))
		}
	}

	return nil, nil
}

func trim(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func contains(s string, e []*ast.Comment) bool {
	for _, b := range e {
		if strings.Contains(strings.ToLower(b.Text), strings.ToLower(s)) {
			return true
		}
	}
	return false
}