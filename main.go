package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"
	"unicode/utf8"

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
	enc := args.Val("enc", os.Getenv("ID1_ENC"))
	sync := args.Has("sync") || os.Getenv("ID1_SYNC") == "true"
	envOp := args.Val("env", "none")

	if opFunc, ok := envOpFunc[envOp]; ok {
		opFunc(args.KeyVal(envOp, "", ""))
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

	if args.Has("connect") {

		cmdIn, cmdOut, disconnect := connect(c)
		lastPing := time.Now().UnixMilli()
		for {
			if time.Now().UnixMilli()-lastPing > 4000 {
				cmdOut <- id1.Command{Op: id1.Set, Key: id1.KK(id, ".ping"), Args: map[string]string{"ttl": "0"}, Data: fmt.Appendf(nil, "%d", lastPing)}
				lastPing = time.Now().UnixMilli()
			}

			select {
			case <-time.After(time.Second * 5):
			case cmd := <-cmdIn:
				switch enc {
				case "base64":
					d, _ := base64.StdEncoding.DecodeString(string(cmd.Data))
					cmd.Data = d
				default:
				}
				if syncOpFunc, ok := syncOpFunc[cmd.Op]; ok && sync {
					syncOpFunc(cmd)
				}
				if utf8.Valid(cmd.Data) {
					fmt.Printf("%s -> %s\n", cmd.String(), string(cmd.Data))
				} else {
					fmt.Printf("%s -> [%d bytes]\n", cmd.String(), len(cmd.Data))
				}
			case <-disconnect:
				fmt.Println("disconnected")
				return
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

	fmt.Println(man)
	os.Exit(0)
}
