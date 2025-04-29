package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	id1 "github.com/qodex/id1-client-go"
)

func connect(id string, c id1.Id1Client, enc string) (chan id1.Command, chan id1.Command, chan bool) {
	if disconnect, err := c.Connect(); err != nil {
		fmt.Printf("connection error: %s\n", err)
		os.Exit(1)
		return nil, nil, nil
	} else {
		fmt.Println("connected")
		cmdIn := make(chan id1.Command, 32)
		cmdOut := make(chan id1.Command)

		c.AddListener(func(cmd id1.Command) {
			if cmd.Op == id1.Unknown {
				return
			}
			switch enc {
			case "base64":
				data, _ := base64.StdEncoding.DecodeString(string(cmd.Data))
				cmd.Data = data
			}
			cmdIn <- cmd
		}, "")

		go send(cmdOut, c)
		go ping(id, cmdOut)

		return cmdIn, cmdOut, disconnect
	}
}

func ping(id string, cmdOut chan id1.Command) {
	for {
		time.Sleep(time.Second * 5)
		cmdOut <- id1.Command{Op: id1.Get, Key: id1.KK(id, ".ping")}
	}
}

func send(cmdOut chan id1.Command, c id1.Id1Client) {
	for {
		cmd := <-cmdOut
		if cmd.Op != id1.Unknown {
			if err := c.Send(cmd); err != nil {
				log.Println(err)
			}
		}
	}
}
