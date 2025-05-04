package main

import (
	"context"
	"fmt"
	"os"
	"regexp"

	id1 "github.com/qodex/id1-client-go"
)

func watch(args id1Args) {
	re, err := regexp.Compile(args.regexp)
	if args.filter && err != nil {
		fmt.Printf("regexp err %s", err)
		os.Exit(1)
	}
	ctx := context.Background()
	cmdOut := make(chan id1.Command, 64)
	go watchDir(args.dir, cmdOut, ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-cmdOut:
			filterPass := true
			if args.filter && re != nil && !re.MatchString(cmd.String()) {
				filterPass = false
			}
			if filterPass {
				os.Stdout.Write(fmt.Appendf(cmd.Bytes(), "%s", "\n\n"))
			}
		}
	}
}
