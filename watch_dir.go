package main

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	id1 "github.com/qodex/id1-client-go"
)

func watchDir(rootpath string, dirCmd chan id1.Command, ctx context.Context) error {
	if watcher, err := fsnotify.NewWatcher(); err != nil {
		return err
	} else if walkErr := filepath.WalkDir(rootpath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("watch dir %s walk err %s", path, err)
			return nil
		}
		if entry.IsDir() {
			watcher.Add(path)
		}
		return nil
	}); walkErr != nil {
		return walkErr
	} else {
		defer watcher.Close()
		watchLoop(watcher, rootpath, dirCmd, ctx)
	}
	return nil
}

var opMap = map[fsnotify.Op]id1.Op{
	fsnotify.Create: id1.Set,
	fsnotify.Write:  id1.Set,
	fsnotify.Remove: id1.Del,
	fsnotify.Rename: id1.Del,
}

func watchLoop(watcher *fsnotify.Watcher, rootpath string, cmdOut chan id1.Command, ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("watch dir: error, chan events is closed, recover")
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-watcher.Errors:
			log.Printf("watch dir err %s", err)
			if !ok {
				return
			}
		case e, ok := <-watcher.Events:
			if !ok {
				return
			}
			ignore := e.Has(fsnotify.Chmod)
			if ignore {
				continue
			}

			// add new dirs to watcher
			if e.Has(fsnotify.Create) {
				if st, err := os.Stat(e.Name); err == nil && st.IsDir() {
					filepath.WalkDir(e.Name, func(path string, d fs.DirEntry, err error) error {
						if d.IsDir() {
							watcher.Add(path)
						}
						return nil
					})
					continue
				}
			}

			cmd := id1.Command{
				Op:  opMap[e.Op],
				Key: id1.K(strings.Trim(e.Name[len(rootpath):], "/")),
			}

			if cmd.Op == id1.Set {
				if data, err := os.ReadFile(e.Name); err == nil {
					cmd.Op = id1.Set
					cmd.Data = data
				}
			}

			if cmd.Op != id1.Unknown {
				cmdOut <- cmd
			}
		}
	}
}
