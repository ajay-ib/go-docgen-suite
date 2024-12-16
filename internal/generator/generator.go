package generator

import (
	"fmt"
	"log"
	"os"
)

func GenerateMarkdown(content string) {
	file, err := os.Create("Gen-README.md")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Documentation generated successfully.")
}
