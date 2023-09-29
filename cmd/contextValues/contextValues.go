package main

import (
	"context"
	"fmt"
)

func enrichContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "api", ctx)
}

func donRequest(ctx context.Context) {
	value := ctx.Value("api")

	fmt.Println(value)
}

func main() {
	donRequest(enrichContext(context.Background()))
}
