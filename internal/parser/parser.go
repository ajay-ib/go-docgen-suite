package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func ParseFile(filename string) (string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	var contentBuilder strings.Builder
	for _, decl := range node.Decls {
		if fn, isFn := decl.(*ast.FuncDecl); isFn {
			funcName := fn.Name.Name
			comments := extractComments(fn)
			contentBuilder.WriteString(fmt.Sprintf("## Function: %s\n\n%s\n\n", funcName, comments))
		}
	}
	return contentBuilder.String(), nil
}

func extractComments(fn *ast.FuncDecl) string {
	if fn.Doc == nil {
		return ""
	}
	var comments []string
	for _, comment := range fn.Doc.List {
		if strings.Contains(comment.Text, fn.Name.Name) {
			comments = append(comments, comment.Text)
		}
	}
	return strings.Join(comments, "\n")
}
