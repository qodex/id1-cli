package main

import (
	"context"
	"fmt"
	"os"

	id1 "github.com/qodex/id1-client-go"
)

func watch(dir string) {
	ctx := context.Background()
	cmdOut := make(chan id1.Command, 64)
	go watchDir(dir, cmdOut, ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-cmdOut:
			if cmd.Op != id1.Unknown {
				os.Stdout.Write(fmt.Appendf(cmd.Bytes(), "%s", "\n\n"))
			}
		}
	}
}
