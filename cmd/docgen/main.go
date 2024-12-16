package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ajay-ib/go-docgen-suite/internal/generator"
	"github.com/ajay-ib/go-docgen-suite/internal/godoc"
	"github.com/ajay-ib/go-docgen-suite/internal/godocweb"
	"github.com/ajay-ib/go-docgen-suite/internal/parser"
	"github.com/ajay-ib/go-docgen-suite/internal/swaggo"
	"github.com/ajay-ib/go-docgen-suite/internal/traversal"
	"github.com/urfave/cli/v2"
)

//go:embed generate-docs.sh
var generateDocsScript string

func main() {
	app := &cli.App{
		Name:  "docgen-service",
		Usage: "Generate documentation for Go services",
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "Generate documentation",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "path",
						Aliases:  []string{"p"},
						Usage:    "Path to the Go service",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "Output directory for the generated documentation",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					root := c.String("path")
					output := c.String("output")
					return generateDocumentation(root, output)
				},
			},
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "Serve documentation via web interface",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "path",
						Aliases:  []string{"p"},
						Usage:    "Path to the Go service",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					root := c.String("path")
					return godocweb.ServeGodocWeb(root)
				},
			},
			{
				Name:    "install-script",
				Aliases: []string{"i"},
				Usage:   "Install the generate-docs.sh script",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "path",
						Aliases:  []string{"p"},
						Usage:    "Path to install the script",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					targetPath := c.String("path")
					return installScript(targetPath)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func generateDocumentation(root, output string) error {
	var contentBuilder strings.Builder

	err := traversal.TraverseFiles(root, func(path string) error {
		content, err := parser.ParseFile(path)
		if err != nil {
			return err
		}
		contentBuilder.WriteString(content)
		return nil
	})
	if err != nil {
		return err
	}

	generator.GenerateMarkdown(contentBuilder.String(), output)

	if err := godoc.GenerateGodoc(root); err != nil {
		return fmt.Errorf("error generating Godoc: %v", err)
	}

	if err := swaggo.GenerateSwaggoDocs(root); err != nil {
		return fmt.Errorf("error generating Swaggo docs: %v", err)
	}

	return nil
}

func installScript(targetPath string) error {
	targetScriptPath := filepath.Join(targetPath, "generate-docs.sh")

	err := os.WriteFile(targetScriptPath, []byte(generateDocsScript), fs.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing script: %v", err)
	}

	fmt.Println("Script installed successfully at", targetScriptPath)
	return nil
}
