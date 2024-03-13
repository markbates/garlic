package garlic

import (
	"context"
	"fmt"
	"io/fs"
	"sync"

	"github.com/markbates/iox"
	"github.com/markbates/plugins/plugcmd"
)

var _ plugcmd.Commander = &Garlic{}

type Garlic struct {
	Cmd  Commander
	FS   fs.FS
	IO   iox.IO
	Name string

	mu sync.RWMutex
}

func (g *Garlic) Main(ctx context.Context, pwd string, args []string) error {
	if g == nil {
		return fmt.Errorf("garlic is nil")
	}

	if len(g.Name) == 0 {
		return fmt.Errorf("command name is required")
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	local := Local{
		FS:   g.FS,
		IO:   g.IO,
		Name: g.Name,
		Root: pwd,
	}

	if local.Exists() {
		return local.Run(ctx, args)
	}

	if g.Cmd == nil {
		return fmt.Errorf("command is nil")
	}

	cmd := g.Cmd

	if sfs, ok := cmd.(SettableFS); ok {
		sfs.SetFS(g.FS)
	}

	if sio, ok := cmd.(SettableIO); ok {
		sio.SetIO(g.IO)
	}

	return cmd.Main(ctx, pwd, args)
}

func (g *Garlic) PluginName() string {
	return fmt.Sprintf("%T", g)
}
