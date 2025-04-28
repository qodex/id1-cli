package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qodex/ff"
	id1 "github.com/qodex/id1-client-go"
)

var syncOpFunc = map[id1.Op]func(cmd id1.Command){
	id1.Set: func(cmd id1.Command) {
		if cmd.Args["ttl"] == "0" {
			return
		}
		path := filepath.Join(os.Getenv("ID1_WORKDIR"), cmd.Key.String())
		ff.CreatePath(path)
		os.WriteFile(path, cmd.Data, os.ModePerm)
	},
	id1.Add: func(cmd id1.Command) {
		path := filepath.Join(os.Getenv("ID1_WORKDIR"), cmd.Key.String())
		ff.CreatePath(path)
		ff.FileAppend(path, cmd.Data)
	},
	id1.Del: func(cmd id1.Command) {
		path := filepath.Join(os.Getenv("ID1_WORKDIR"), cmd.Key.String())
		if _, err := os.Stat(path); err == nil {
			if err := os.RemoveAll(path); err != nil {
				fmt.Printf("del err %s", err)
			}
		}
	},
	id1.Mov: func(cmd id1.Command) {
		oldpath := filepath.Join(os.Getenv("ID1_WORKDIR"), cmd.Key.String())
		newpath := filepath.Join(os.Getenv("ID1_WORKDIR"), string(cmd.Data))
		ff.CreatePath(newpath)
		if _, err := os.Stat(oldpath); err == nil {
			if err := os.Rename(oldpath, newpath); err != nil {
				fmt.Printf("mov err\nold=%s\nnew=%s\nerr=%s\n", oldpath, newpath, err)
			}
		}
	},
}
