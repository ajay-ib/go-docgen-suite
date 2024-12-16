package swaggo

import (
	"os/exec"
)

func GenerateSwaggoDocs(root string) error {
	cmd := exec.Command("swag", "init", "--dir", root)
	return cmd.Run()
}
