package godoc

import (
	"os/exec"
)

func GenerateGodoc(root string) error {
	cmd := exec.Command("godoc", "-goroot", root)
	return cmd.Run()
}
