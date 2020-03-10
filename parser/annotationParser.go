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

// Return a list of all functions whose annotations match with the test names in the config
// @param test Initialized test struct from github.com/evoila/infraTESTure/config
// @param path Path to the test project
// @return []string Name list of available functions in the test project
func GetFunctionNames(test config.Test, path string) ([]string, error) {
	var methodNames []string

	fileSet := token.NewFileSet() // positions are relative to fileSet
	file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Iterate through all function declarations to check if any annotation matches with the
	// tests provided in the configuration.yml
	for _, decl := range file.Decls {
		if fun, ok := decl.(*ast.FuncDecl); ok {
			if fun.Doc != nil && contains("@"+test.Name, fun.Doc.List) {
				methodNames = append(methodNames, fun.Name.Name)
			}
		}
	}

	return methodNames, nil
}

// Return a list of all annotations in a given go file
// @param path Path to the test project
// @return []string List of annotations from functions within the test project
func GetAnnotations(path string) []string {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of all comments in the file
	comments := ast.NewCommentMap(fileSet, file, file.Comments).Comments()

	var testNames []string

	// Iterate comments and check if its an annotation
	for _, comment := range comments {
		if strings.HasPrefix(comment.Text(), "@") {
			testNames = append(testNames, trim(comment.Text()))
		}
	}

	return testNames
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