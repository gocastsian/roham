package job

import "context"

func SayHelloInPersian(ctx context.Context, name string) (string, error) {
	return "سلام " + name, nil
}
