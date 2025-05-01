package main

import (
	"fmt"
	"os"

	id1 "github.com/qodex/id1-client-go"
)

func cmdExec(cmd id1.Command, c id1.Id1Client) {
	if data, err := c.Exec(cmd); err != nil {
		fmt.Println(err)
	} else {
		os.Stdout.Write(data)
	}
}
