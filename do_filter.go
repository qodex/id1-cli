package main

import (
	"fmt"
	"os"
	"regexp"

	id1 "github.com/qodex/id1-client-go"
)

func filter(args id1Args) {
	regexp, err := regexp.Compile(args.regexp)
	if err != nil {
		fmt.Printf("regexp err %s", err)
		os.Exit(1)
	}
	cmdIn := make(chan id1.Command, 32)
	eof := make(chan bool)
	go scanCommands(cmdIn, eof)
	for {
		select {
		case cmd := <-cmdIn:
			if cmd.Op != id1.Unknown && regexp.MatchString(cmd.String()) {
				os.Stdout.Write(fmt.Appendf(nil, "%s\n\n", string(cmd.Bytes())))
			}
		case <-eof:
			os.Exit(0)
		}
	}
}
