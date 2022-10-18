package garlic

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/markbates/iox"
	"github.com/markbates/plugins"
	"github.com/stretchr/testify/require"
)

type commander func(ctx context.Context, pwd string, args []string) error

func (c commander) PluginName() string {
	return plugins.Name(c)
}

func (c commander) Main(ctx context.Context, pwd string, args []string) error {
	return c(ctx, pwd, args)
}

func Test_Garlic_Main_NoLocal(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	g := &Garlic{
		Name: "commander",
	}

	fn := func(ctx context.Context, pwd string, args []string) error {
		g.Exit = 42

		return nil
	}

	g.Cmd = commander(fn)

	err := g.Main(context.Background(), ".", []string{})
	r.NoError(err)
	r.Equal(42, g.Exit)
}

func Test_Garlic_Main_Local(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	oi := iox.Buffer{}

	g := &Garlic{
		Name: "commander",
		IO:   oi.IO(),
	}

	fn := func(ctx context.Context, pwd string, args []string) error {
		return fmt.Errorf("this should not be called")
	}

	g.Cmd = commander(fn)

	pwd := "testdata"
	err := g.Main(context.Background(), pwd, []string{"foo", "bar"})
	r.Error(err)

	act := oi.Out.String()
	act = strings.TrimSpace(act)
	exp := `[foo bar]`

	r.Equal(exp, act)

	r.Equal(1, g.Exit)
}
