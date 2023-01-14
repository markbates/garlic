package garlic

import (
	"context"
	"fmt"
	"os"
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

func Test_Garlic_Local(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	oi := iox.Buffer{}
	var fn commander = func(ctx context.Context, pwd string, args []string) error {
		fmt.Fprintln(&oi.Out, "should not be called")
		return nil
	}

	g := &Garlic{
		Name: "test",
		Cmd:  fn,
		IO:   oi.IO(),
		FS:   os.DirFS("."),
	}

	err := g.Main(context.Background(), "testdata/local/plugs", []string{"hello"})
	fmt.Println(oi.Out.String())
	r.NoError(err)

	act := oi.Out.String()
	r.Contains(act, "args: [hello]")
}

func Test_Garlic_NoLocal(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	oi := iox.Buffer{}
	var fn commander = func(ctx context.Context, pwd string, args []string) error {
		fmt.Fprintln(&oi.Out, "should be called")
		return nil
	}

	g := &Garlic{
		Name: "test",
		Cmd:  fn,
		IO:   oi.IO(),
		FS:   os.DirFS("."),
	}

	err := g.Main(context.Background(), "testdata/local/noplugs", []string{"hello"})
	fmt.Println(oi.Out.String())
	r.NoError(err)

	act := oi.Out.String()
	r.Contains(act, "should be called")
}
