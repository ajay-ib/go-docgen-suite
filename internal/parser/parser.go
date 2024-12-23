package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type ParsedContent struct {
	DeveloperDoc string
	CopilotDoc   string
}

func ParseFile(filename string) (ParsedContent, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return ParsedContent{}, err
	}

	var developerDoc, copilotDoc strings.Builder
	for _, decl := range node.Decls {
		if fn, isFn := decl.(*ast.FuncDecl); isFn {
			funcName := fn.Name.Name
			comments := extractComments(fn)
			developerDoc.WriteString(fmt.Sprintf("## Function: %s\n\n%s\n\n", funcName, comments))
			copilotDoc.WriteString(fmt.Sprintf(`{"type": "function", "name": "%s", "description": "%s"}`, funcName, comments))
		}
	}
	return ParsedContent{
		DeveloperDoc: developerDoc.String(),
		CopilotDoc:   copilotDoc.String(),
	}, nil
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
