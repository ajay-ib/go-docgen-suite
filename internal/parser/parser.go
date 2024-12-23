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
    developerDoc.WriteString(fmt.Sprintf("Signature: %s\n", formatFuncSignature(decl)))
    copilotDoc.WriteString(fmt.Sprintf(`{"type": "function", "name": "%s", "description": "%s", "signature": "%s"}`, decl.Name.Name, decl.Doc.Text(), formatFuncSignature(decl)))
}

func formatFuncSignature(decl *ast.FuncDecl) string {
    var params, results []string
    for _, param := range decl.Type.Params.List {
        paramType := fmt.Sprintf("%s", param.Type)
        for _, name := range param.Names {
            params = append(params, fmt.Sprintf("%s %s", name.Name, paramType))
        }
    }
    if decl.Type.Results != nil {
        for _, result := range decl.Type.Results.List {
            resultType := fmt.Sprintf("%s", result.Type)
            if len(result.Names) > 0 {
                for _, name := range result.Names {
                    results = append(results, fmt.Sprintf("%s %s", name.Name, resultType))
                }
            } else {
                results = append(results, resultType)
            }
        }
    }
    return fmt.Sprintf("func(%s) (%s)", strings.Join(params, ", "), strings.Join(results, ", "))
}