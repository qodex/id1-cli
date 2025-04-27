package main

import (
	"fmt"
	"log"
	"os"

	id1 "github.com/qodex/id1-client-go"
)

func connect(c id1.Id1Client) (chan id1.Command, chan id1.Command, chan bool) {
	if disconnect, err := c.Connect(); err != nil {
		fmt.Printf("connection error: %s\n", err)
		os.Exit(1)
		return nil, nil, nil
	} else {
		fmt.Println("connected")
		cmdIn := make(chan id1.Command, 32)
		cmdOut := make(chan id1.Command)
		c.AddListener(func(cmd id1.Command) {
			if cmd.Op != id1.Unknown {
				cmdIn <- cmd
			}
		}, "")
		go func() {
			for {
				cmd := <-cmdOut
				if err := c.Send(cmd); err != nil {
					log.Println(err)
				}
			}
		}()
		go scanCommands(cmdOut)
		return cmdIn, cmdOut, disconnect
	}
}
