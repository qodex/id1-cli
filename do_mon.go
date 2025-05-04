package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"time"

	id1 "github.com/qodex/id1-client-go"
)

func mon(args id1Args, c id1.Id1Client) {
	re, err := regexp.Compile(args.regexp)
	if args.filter && err != nil {
		fmt.Printf("regexp err %s", err)
		os.Exit(1)
	}
	cmdIn, cmdOut, disconnect := connect(args.id, c, args.enc)
	eof := make(chan bool)
	go scanCommands(cmdOut, eof)
	for {
		select {
		case cmd := <-cmdIn:
			filterPass := true
			if args.filter && re != nil && !re.MatchString(cmd.String()) {
				filterPass = false
			}
			if filterPass {
				os.Stdout.Write(fmt.Appendf(nil, "%s\n\n", string(cmd.Bytes())))
			}
		case <-eof:
			os.Exit(0)
		case <-disconnect:
			fmt.Println("disconnected")
			os.Exit(0)
		}
	}
}

func connect(id string, c id1.Id1Client, enc string) (chan id1.Command, chan id1.Command, chan bool) {
	if disconnect, err := c.Connect(); err != nil {
		fmt.Printf("connection error: %s\n", err)
		os.Exit(1)
		return nil, nil, nil
	} else {
		fmt.Printf("connected\n\n")
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
				fmt.Printf("send err %s\n", err)
			}
		}
	}
}
