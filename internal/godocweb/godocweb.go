package godocweb

import (
	"os/exec"
)

func ServeGodocWeb(root string) error {
	cmd := exec.Command("godoc", "-http=:6060", "-goroot", root)
	return cmd.Start()
}
