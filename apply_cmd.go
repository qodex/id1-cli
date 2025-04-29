package main

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/qodex/ff"
	id1 "github.com/qodex/id1-client-go"
)

func applyCmd(cmd id1.Command, workdir string) {
	if f, ok := syncOpFunc[cmd.Op]; ok {
		f(cmd, workdir)

		if utf8.Valid(cmd.Data) {
			fmt.Printf("%s\n%s\n\n", cmd.String(), string(cmd.Data))
		} else {
			fmt.Printf("%s\n[%d bytes]\n\n", cmd.String(), len(cmd.Data))
		}
	}
}

var syncOpFunc = map[id1.Op]func(cmd id1.Command, workdir string){
	id1.Set: func(cmd id1.Command, workdir string) {
		if cmd.Args["ttl"] == "0" {
			return
		}
		path := filepath.Join(workdir, cmd.Key.String())
		ff.CreatePath(path)
		os.WriteFile(path, cmd.Data, os.ModePerm)
	},
	id1.Add: func(cmd id1.Command, workdir string) {
		path := filepath.Join(workdir, cmd.Key.String())
		ff.CreatePath(path)
		ff.FileAppend(path, cmd.Data)
	},
	id1.Del: func(cmd id1.Command, workdir string) {
		path := filepath.Join(workdir, cmd.Key.String())
		if _, err := os.Stat(path); err == nil {
			if err := os.RemoveAll(path); err != nil {
				fmt.Printf("del err %s", err)
			}
		}
	},
	id1.Mov: func(cmd id1.Command, workdir string) {
		oldpath := filepath.Join(workdir, cmd.Key.String())
		newpath := filepath.Join(workdir, string(cmd.Data))
		ff.CreatePath(newpath)
		if _, err := os.Stat(oldpath); err == nil {
			if err := os.Rename(oldpath, newpath); err != nil {
				fmt.Printf("mov err\nold=%s\nnew=%s\nerr=%s\n", oldpath, newpath, err)
			}
		}
	},
}
