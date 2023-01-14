package garlic

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/markbates/iox"
)

type Local struct {
	FS   fs.FS
	Root string
	Name string
	IO   iox.IO
}

func (l Local) Run(ctx context.Context, args []string) error {
	cmd := exec.CommandContext(ctx, "go", append([]string{"run", "."}, args...)...)
	cmd.Dir = filepath.Join(l.Root, "garlic", l.Name)
	cmd.Stdout = l.IO.Stdout()
	cmd.Stderr = l.IO.Stderr()
	cmd.Stdin = l.IO.Stdin()

	return cmd.Run()
}

func (l Local) String() string {
	x := struct {
		Root   string
		CmdDir string
	}{
		Root:   l.Root,
		CmdDir: filepath.Join(l.Root, "garlic", l.Name),
	}

	return fmt.Sprintf("%+v", x)
}

func (l Local) Exists() bool {
	if l.FS == nil {
		return false
	}

	pwd := l.Root
	if len(pwd) == 0 {
		pwd = "."
	}

	cdir := filepath.Join(l.Root, "garlic", l.Name)
	if len(cdir) == 0 {
		cdir = filepath.Join(pwd, "cmd", "garlic")
	}

	if filepath.IsAbs(pwd) {
		if _, err := os.Stat(pwd); err != nil {
			return false
		}

		if _, err := os.Stat(cdir); err != nil {
			return false
		}
		return true
	}

	if _, err := fs.Stat(l.FS, pwd); err != nil {
		return false
	}

	if _, err := fs.Stat(l.FS, cdir); err != nil {
		return false
	}

	return true
}
