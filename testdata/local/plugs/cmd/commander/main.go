package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/markbates/garlic/testdata/local/plugs/cmd/commander/cli"
	"github.com/markbates/iox"
	"github.com/markbates/plugins"
)

func main() {
	args := os.Args[1:]

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cmd := cli.New(pwd)

	plugs := cmd.Plugins()
	plugs = append(plugs, &Sub{})

	// your plugins here:
	// cmd.Plugins = append(cmd.Plugins, ...)
	fn := func() plugins.Plugins {
		return plugs
	}

	cmd.Feeder = fn

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ctx, cancel = signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	fmt.Println("starting")
	err = cmd.Main(ctx, pwd, args)
	fmt.Println("err:", err)
	fmt.Println("args:", cmd.Result.Args)

	if err != nil {
		cmd.Exit(1, err)
	}

}

var _ iox.IOSetable = &Sub{}

type Sub struct {
	iox.IO
}

func (s *Sub) SetStdio(oi iox.IO) {
	s.IO = oi
}

func (s Sub) PluginName() string {
	return plugins.Name(s)
}

func (s Sub) CmdName() string {
	return "sub"
}

func (s Sub) Main(ctx context.Context, pwd string, args []string) error {
	fmt.Fprintln(s.Stdout(), "Hello from Sub!")
	return nil
}

func (s Sub) MarkedPlug() {}

func (s Sub) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"name":  s.PluginName(),
		"cmd":   s.CmdName(),
		"stdio": s.IO,
	})
}

func (s Sub) String() string {
	b, _ := s.MarshalJSON()
	return string(b)
}
