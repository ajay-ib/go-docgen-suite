package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ajay-ib/go-docgen-suite/internal/generator"
	"github.com/ajay-ib/go-docgen-suite/internal/godoc"
	"github.com/ajay-ib/go-docgen-suite/internal/godocweb"
	"github.com/ajay-ib/go-docgen-suite/internal/parser"
	"github.com/ajay-ib/go-docgen-suite/internal/swaggo"
	"github.com/ajay-ib/go-docgen-suite/internal/traversal"
	"github.com/urfave/cli/v2"
)

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
