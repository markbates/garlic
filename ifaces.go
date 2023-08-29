package garlic

import "context"

type Commander interface {
	Main(ctx context.Context, pwd string, args []string) error
}
