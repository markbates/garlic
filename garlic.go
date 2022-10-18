package garlic

import (
	"context"
	"errors"
	"fmt"
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

	if _, err := os.Stat(main); err != nil {
		return c.Main(ctx, pwd, args)
	}

	bargs := []string{"run", "-v", fmt.Sprintf("./cmd/%s/%s", g.Name, "main.go")}
	bargs = append(bargs, args...)

	cmd := exec.CommandContext(ctx, "go", bargs...)
	cmd.Dir = pwd
	cmd.Env = os.Environ()
	cmd.Stdin = oi.Stdin()
	cmd.Stdout = oi.Stdout()
	cmd.Stderr = oi.Stderr()

	if err := cmd.Run(); err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			g.Exit = e.ExitCode()
		}

		return fmt.Errorf("failed to run %q in %q: %w", cmd.Args, pwd, err)
	}

	return nil
}

func (g Garlic) PluginName() string {
	return plugins.Name(g)
}
