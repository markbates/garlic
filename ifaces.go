package garlic

import (
	"context"
	"io/fs"

	"github.com/markbates/iox"
)

type Commander interface {
	Main(ctx context.Context, pwd string, args []string) error
}

type SettableIO interface {
	SetIO(io iox.IO)
}

type SettableFS interface {
	SetFS(fs fs.FS)
}
