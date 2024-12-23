package parser

import (
    "encoding/json"
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

type CopilotEntry struct {
    Type        string `json:"type"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Signature   string `json:"signature"`
}

func ParseFile(filename string) (ParsedContent, error) {
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
    if err != nil {
        return ParsedContent{}, err
    }

    var developerDoc strings.Builder
    var copilotEntries []CopilotEntry

    for _, decl := range node.Decls {
        switch d := decl.(type) {
        case *ast.GenDecl:
            for _, spec := range d.Specs {
                switch s := spec.(type) {
                case *ast.TypeSpec:
                    processTypeSpec(s, d.Doc, &developerDoc, &copilotEntries)
                }
            }
        case *ast.FuncDecl:
            processFuncDecl(d, &developerDoc, &copilotEntries)
        }
    }

    copilotDoc, err := json.MarshalIndent(copilotEntries, "", "  ")
    if err != nil {
        return ParsedContent{}, err
    }

    return ParsedContent{
        DeveloperDoc: developerDoc.String(),
        CopilotDoc:   string(copilotDoc),
    }, nil
}

func processTypeSpec(spec *ast.TypeSpec, doc *ast.CommentGroup, developerDoc *strings.Builder, copilotEntries *[]CopilotEntry) {
    description := ""
    if doc != nil {
        description = doc.Text()
        developerDoc.WriteString(description + "\n")
    }
    developerDoc.WriteString(fmt.Sprintf("Type: %s\n", spec.Name.Name))
    *copilotEntries = append(*copilotEntries, CopilotEntry{
        Type:        "type",
        Name:        spec.Name.Name,
        Description: description,
        Signature:   "",
    })
}

func processFuncDecl(decl *ast.FuncDecl, developerDoc *strings.Builder, copilotEntries *[]CopilotEntry) {
    description := ""
    if decl.Doc != nil {
        description = decl.Doc.Text()
        developerDoc.WriteString(description + "\n")
    }
    developerDoc.WriteString(fmt.Sprintf("Function: %s\n", decl.Name.Name))
    developerDoc.WriteString(fmt.Sprintf("Signature: %s\n", formatFuncSignature(decl)))
    *copilotEntries = append(*copilotEntries, CopilotEntry{
        Type:        "function",
        Name:        decl.Name.Name,
        Description: description,
        Signature:   formatFuncSignature(decl),
    })
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