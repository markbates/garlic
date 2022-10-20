package garlic

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/markbates/plugins/plugcmd"
)

type Garlic struct {
	Cmd  plugcmd.Commander
	Name string
	IO   iox.IO
	Exit int
}

func (g *Garlic) Main(ctx context.Context, pwd string, args []string) error {
	if g == nil {
		return fmt.Errorf("garlic is nil")
	}

	if len(g.Name) == 0 {
		return fmt.Errorf("command name is required")
	}

	c := g.Cmd

	oi := g.IO

	main := filepath.Join(pwd, "cmd", g.Name, "main.go")

	fn := c.Main

	if _, err := os.Stat(main); err == nil {
		fn = g.Local
	}

	err := fn(ctx, pwd, args)
	if err == nil {
		return nil
	}

	b, be := io.ReadAll(oi.Stderr())
	if be != nil {
		return be
	}

	err = fmt.Errorf("%w: %s", err, string(b))
	var e *exec.ExitError
	if errors.As(err, &e) {
		g.Exit = e.ExitCode()
		return err
	}

	var ex ExitError
	if errors.As(err, &ex) {
		g.Exit = ex.Code
		return err
	}

	return err

}

func (g *Garlic) Local(ctx context.Context, pwd string, args []string) error {
	oi := g.IO

	// path to local folder containing main.go
	cdir := filepath.Join(pwd, "cmd", g.Name)

	bargs := []string{"build", "-v"}

	// create a temp dir to put the binary
	odir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}

	// clean up the temp dir
	defer os.RemoveAll(odir)

	// path to the binary
	bin := filepath.Join(odir, g.Name)

	// add the output path to the build args
	bargs = append(bargs, "-o", bin)

	// add the path to the local folder to build
	bargs = append(bargs, "./"+cdir)

	cmd := exec.CommandContext(ctx, "go", bargs...)
	cmd.Env = os.Environ()
	cmd.Stdin = oi.Stdin()
	cmd.Stdout = oi.Stdout()
	cmd.Stderr = oi.Stderr()

	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.CommandContext(ctx, bin, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = oi.Stdin()
	cmd.Stdout = oi.Stdout()
	cmd.Stderr = oi.Stderr()

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (g Garlic) PluginName() string {
	return plugins.Name(g)
}
