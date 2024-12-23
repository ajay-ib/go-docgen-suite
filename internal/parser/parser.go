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
        switch d := decl.(type) {
        case *ast.GenDecl:
            for _, spec := range d.Specs {
                switch s := spec.(type) {
                case *ast.TypeSpec:
                    processTypeSpec(s, d.Doc, &developerDoc, &copilotDoc)
                }
            }
        case *ast.FuncDecl:
            processFuncDecl(d, &developerDoc, &copilotDoc)
        }
    }
    return ParsedContent{
        DeveloperDoc: developerDoc.String(),
        CopilotDoc:   copilotDoc.String(),
    }, nil
}

func processTypeSpec(spec *ast.TypeSpec, doc *ast.CommentGroup, developerDoc, copilotDoc *strings.Builder) {
    if doc != nil {
        developerDoc.WriteString(doc.Text() + "\n")
    }
    developerDoc.WriteString(fmt.Sprintf("Type: %s\n", spec.Name.Name))
    copilotDoc.WriteString(fmt.Sprintf(`{"type": "type", "name": "%s", "description": "%s"}`, spec.Name.Name, doc.Text()))
}

func processFuncDecl(decl *ast.FuncDecl, developerDoc, copilotDoc *strings.Builder) {
    if decl.Doc != nil {
        developerDoc.WriteString(decl.Doc.Text() + "\n")
    }
    developerDoc.WriteString(fmt.Sprintf("Function: %s\n", decl.Name.Name))
    copilotDoc.WriteString(fmt.Sprintf(`{"type": "function", "name": "%s", "description": "%s"}`, decl.Name.Name, decl.Doc.Text()))
}