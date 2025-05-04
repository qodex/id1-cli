package main

import (
	id1 "github.com/qodex/id1-client-go"
	"os"
)

func cmdExec(cmd id1.Command, c id1.Id1Client) {
	if data, err := c.Exec(cmd); err != nil {
		os.Exit(1)
	} else {
		os.Stdout.Write(data)
	}
}
