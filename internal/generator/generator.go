package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GenerateMarkdown(content, outputDir string) {
	outputPath := filepath.Join(outputDir, "Gen-README.md")
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Documentation generated successfully at", outputPath)
}
