package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/qodex/ff"
	id1 "github.com/qodex/id1-client-go"
)

func main() {
	godotenv.Load()
	args := ff.NewOsArgs(os.Args)

	url := args.Val("url", args.WithPrefix("http", os.Getenv("ID1_URL")))
	id := args.Val("id", os.Getenv("ID1_ID"))
	keyPath := args.Val("key", os.Getenv("ID1_KEY"))
	encode := args.Has("base64") || args.Has("b64") || args.Has("enc")

	if args.Has("env") && args.Has("set") {
		kv := strings.Split(args.RestAfter("set", ""), "=")
		if len(kv) > 0 {
			envVar := strings.ToUpper("ID1_" + kv[0])
			val := strings.Join(kv[1:], " ")
			ff.NewFsProps(".env").Set(envVar, val)
		}
		os.Exit(0)
	}

	if args.Has("env") && args.Has("get") {
		kv := strings.Split(args.RestAfter("get", ""), "=")
		if len(kv) > 0 {
			envVar := strings.ToUpper("ID1_" + kv[0])
			fmt.Println(os.Getenv(envVar))
		}
		os.Exit(0)
	}

	c, err := id1.NewHttpClient(url)
	if err != nil {
		fmt.Printf("unvalid url: %s", err)
		os.Exit(1)
	}

	if encode {
		prep := []func(cmd *id1.Command) error{
			func(cmd *id1.Command) error {
				if len(cmd.Data) > 0 {
					cmd.Data = []byte(base64.StdEncoding.EncodeToString(cmd.Data))
				}
				return nil
			},
		}
		postp := []func(data []byte, err error) ([]byte, error){}
		proxy := id1.NewId1ClientProxy(c, prep, postp)
		c = proxy
	}

	if create := args.Val("create", ""); len(create) > 0 {
		pubKey := id1.KK(create, "pub", "key")
		if _, err := c.Get(pubKey); err == nil {
			fmt.Println("id already exists")
			os.Exit(1)
		}

		publicKey := ""
		privateKey := ""
		if keyPair, err := genKey(1024); err != nil {
			fmt.Println("key gen error", err)
			os.Exit(1)
		} else {
			publicKey = keyPair.public
			privateKey = keyPair.private
		}
		if err := c.Set(pubKey, []byte(publicKey)); err != nil {
			fmt.Printf("set %s error %s\n", pubKey, err)
			os.Exit(1)
		} else {
			fmt.Print(privateKey)
			os.Exit(0)
		}
	}

	if key, err := ff.ReadString(keyPath); err != nil || len(key) == 0 {
		fmt.Printf("error reading key %s", err)
		os.Exit(1)
	} else if err := c.Authenticate(id, key); err != nil {
		fmt.Printf("authentication failed %s %d %s\n", id, len(key), url)
		os.Exit(1)
	}

	if args.Has("connect") {
		cmdIn, cmdOut, disconnect := connect(c)
		for {
			select {
			case <-disconnect:
				fmt.Println("disconnected")
				return
			case <-time.After(time.Second * 10):
				cmdOut <- id1.Command{Op: id1.Set, Key: id1.KK(id, "ping"), Data: fmt.Appendf(nil, "%d", time.Now().UnixMilli())}
			case cmd := <-cmdIn:
				if encode {
					d, _ := base64.StdEncoding.DecodeString(string(cmd.Data))
					cmd.Data = d
				}
				if cmd.Key.Name != "ping" {
					fmt.Println(cmd)
				}
			}
		}

	}

	cmdStr := args.Find(func(arg string) bool {
		cmd, err := id1.ParseCommand([]byte(arg))
		return err == nil && cmd.Op != id1.Unknown
	}, "")

	if len(cmdStr) > 0 {
		cmdData := args.RestAfter(cmdStr, "")
		cmd, _ := id1.ParseCommand(fmt.Appendf(nil, "%s\n%s", cmdStr, cmdData))
		if data, err := c.Exec(cmd); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			os.Stdout.Write(data)
		}
		os.Exit(0)
	}

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)
	<-ctrlC
}
