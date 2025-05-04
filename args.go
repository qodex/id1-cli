package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/qodex/ff"
	id1 "github.com/qodex/id1-client-go"
)

type id1Args struct {
	args                                                           ff.OsArgs
	env, create, watch, mon, apply, serve, filter                  bool
	dir, url, id, createId, keyPath, key, enc, envOp, port, regexp string
	cmd                                                            *id1.Command
}

func (t id1Args) KeyVal(name, key, value string) (string, string) {
	return t.args.KeyVal(name, key, value)
}

func getArgs() id1Args {
	godotenv.Load()
	args := ff.NewOsArgs(os.Args)
	id1Args := id1Args{
		args: args,
	}
	id1Args.dir = args.Val("dir", os.Getenv("ID1_DIR"))
	if wd, err := filepath.Abs(id1Args.dir); err == nil {
		id1Args.dir = wd
	} else if wd, err := os.Getwd(); err == nil {
		id1Args.dir = wd
	}
	id1Args.url = args.Val("url", args.WithPrefix("http", os.Getenv("ID1_URL")))
	id1Args.id = args.Val("id", os.Getenv("ID1_ID"))
	id1Args.create = args.Has("create")
	id1Args.createId = args.Val("create", "")
	id1Args.keyPath = args.Val("key", os.Getenv("ID1_KEY"))
	id1Args.enc = args.Val("enc", os.Getenv("ID1_ENC"))
	id1Args.env = args.Has("env")
	id1Args.envOp = args.Val("env", "")
	id1Args.filter = args.Has("filter")
	id1Args.regexp = args.Val("filter", ".")
	id1Args.watch = args.Has("watch")
	id1Args.mon = args.Has("mon")
	id1Args.apply = args.Has("apply")
	id1Args.serve = args.Has("serve")
	id1Args.port = args.Val("serve", args.Val("p", "8080"))

	if str, err := ff.ReadString(id1Args.keyPath); len(id1Args.keyPath) > 0 && err != nil {
		fmt.Printf("error reading key %s", err)
	} else {
		id1Args.key = str
	}

	cmdStr := args.Find(func(arg string) bool {
		cmd, err := id1.ParseCommand([]byte(arg))
		return err == nil && cmd.Op != id1.Unknown
	}, "")

	if len(cmdStr) > 0 {
		cmdData := args.RestAfter(cmdStr, "")
		stdinData := ff.ScanStdinBytes()
		if len(stdinData) > 0 {
			cmdData = string(stdinData)
		}
		if cmd, err := id1.ParseCommand(fmt.Appendf(nil, "%s\n%s", cmdStr, cmdData)); err == nil {
			id1Args.cmd = &cmd
		}
	}

	return id1Args
}
