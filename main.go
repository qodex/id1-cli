package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/qodex/ff"
	id1 "github.com/qodex/id1-client-go"
)

func main() {
	godotenv.Load()
	args := ff.NewOsArgs(os.Args)

	dir := args.Val("dir", os.Getenv("ID1_DIR"))
	if wd, err := filepath.Abs(dir); err == nil {
		dir = wd
	} else if wd, err := os.Getwd(); err == nil {
		dir = wd
	}
	url := args.Val("url", args.WithPrefix("http", os.Getenv("ID1_URL")))
	id := args.Val("id", os.Getenv("ID1_ID"))
	createId := args.Val("create", "")
	keyPath := args.Val("key", os.Getenv("ID1_KEY"))
	enc := args.Val("enc", os.Getenv("ID1_ENC"))
	env := args.Has("env")
	envOp := args.Val("env", "")
	watch := args.Has("watch")
	mon := args.Has("mon")
	apply := args.Has("apply")

	if env {
		if f, ok := envOpFunc[envOp]; ok {
			f(args.KeyVal(envOp, "", ""))
		}
	}

	c, err := id1.NewHttpClient(url)
	if err != nil {
		fmt.Printf("unvalid url: %s", err)
		os.Exit(1)
	}

	switch enc {
	case "base64":
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
	default:
	}

	if len(createId) > 0 {
		pubKey := id1.KK(createId, "pub", "key")
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

	var key string
	if str, err := ff.ReadString(keyPath); len(keyPath) > 0 && err != nil {
		fmt.Printf("error reading key %s", err)
	} else {
		key = str
	}

	if len(id) > 0 && len(key) > 0 {
		if err := c.Authenticate(id, key); err != nil {
			fmt.Printf("authentication failed %s %d %s\n", id, len(key), url)
			os.Exit(1)
		}
	}

	if mon {
		cmdIn, cmdOut, disconnect := connect(id, c, enc)
		go scanCommands(cmdOut)
		for {
			select {
			case cmd := <-cmdIn:
				os.Stdout.Write(fmt.Appendf(nil, "%s\n\n", string(cmd.Bytes())))
			case <-disconnect:
				fmt.Println("disconnected")
				os.Exit(0)
			}
		}
	}

	if apply {
		cmdIn := make(chan id1.Command)
		ctrlC := make(chan os.Signal, 1)
		signal.Notify(ctrlC, os.Interrupt)
		go scanCommands(cmdIn)
		for {
			select {
			case cmd := <-cmdIn:
				applyCmd(cmd, dir)
			case <-ctrlC:
				os.Exit(0)
			}
		}
	}

	if watch {
		ctx := context.Background()
		dirCmd := make(chan id1.Command, 64)
		go watchDir(dir, dirCmd, ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case cmd := <-dirCmd:
				if cmd.Op != id1.Unknown {
					os.Stdout.Write(fmt.Appendf(cmd.Bytes(), "%s", "\n\n"))
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
		stdinData := scanData()
		if len(stdinData) > 0 {
			cmdData = string(stdinData)
		}
		cmd, _ := id1.ParseCommand(fmt.Appendf(nil, "%s\n%s", cmdStr, cmdData))
		if data, err := c.Exec(cmd); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			os.Stdout.Write(data)
		}
		os.Exit(0)
	}

	fmt.Println(man)
	os.Exit(0)
}
