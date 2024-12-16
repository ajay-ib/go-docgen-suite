package godoc

import (
	"os/exec"
	"path/filepath"
)

func GenerateGodoc(root string) error {
	godocPath := filepath.Join("vendor", "golang.org", "x", "tools", "cmd", "godoc")
	cmd := exec.Command(godocPath, "-goroot", root)
	return cmd.Run()
}
