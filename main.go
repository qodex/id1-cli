package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	id1 "github.com/qodex/id1-client-go"
)

func main() {
	url := flag.String("url", os.Getenv("ID1_URL"), "id1 API endpoint url")
	id := flag.String("id", os.Getenv("ID1_ID"), "id1 id")
	create := flag.String("create", "", "id to create")
	keyPath := flag.String("key", os.Getenv("ID1_KEY_PATH"), "path to a key PEM file")
	flag.Parse()
	cmd, _ := parseCommand(strings.Join(flag.Args(), " "), scanData())

	if len(*id) == 0 && len(*create) == 0 && !(cmd.Key.Pub && cmd.Key.Name == "key") {
		fmt.Println(man)
		os.Exit(1)
	}

	if len(*id) > 0 && len(*keyPath) == 0 {
		fmt.Println(man)
		os.Exit(1)
	}

	c, err := id1.NewHttpClient(*url)
	if err != nil {
		fmt.Printf("unvalid url: %s", err)
		os.Exit(1)
	}

	if len(*create) > 0 {
		if privateKey, err := createId(*create, c); err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		} else {
			fmt.Print(privateKey)
			os.Exit(0)
		}
	}

	key, err := readKey(*keyPath)
	if !(cmd.Key.Pub && cmd.Key.Name == "key") && err != nil {
		fmt.Printf("error reading key: %s", err)
		os.Exit(1)
	}

	if err := c.Authenticate(*id, key); len(key) > 0 && err != nil {
		fmt.Println("Authentication failed")
		os.Exit(1)
	}

	if len(cmd.Key.String()) > 0 {
		// command provided, execute
		if data, err := c.Exec(cmd); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		} else {
			os.Stdout.Write(data)
		}
		os.Exit(0)
	} else if disconnect, err := c.Connect(); err != nil { // command not provided, connect
		fmt.Printf("Connection error: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("connected to %s\n\n", *url)
		c.AddListener(func(cmd id1.Command) {
			fmt.Println(cmd.String())
		}, "")
		cmdOut := make(chan id1.Command)
		go scanCommands(cmdOut)
		go func() {
			for {
				c.Send(<-cmdOut)
			}
		}()
		go func() {
			for {
				<-disconnect
				fmt.Println("disconnected")
				os.Exit(0)
			}
		}()
	}

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)
	<-ctrlC
}
