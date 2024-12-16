package traversal

import (
	"fmt"
	"os"
	"path/filepath"
)

func TraverseFiles(root string, processFile func(string) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %q: %v", path, err)
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			fmt.Println("Processing file:", path)
			if err := processFile(path); err != nil {
				return fmt.Errorf("error processing file %q: %v", path, err)
			}
		}
		return nil
	})
}
