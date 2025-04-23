package main

import (
	"os"
	"path/filepath"
	"strings"
)

func readKey(path string) (string, error) {
	if !strings.HasPrefix(path, "/") || !strings.HasPrefix(path, "~") {
		wd, _ := os.Getwd()
		path = filepath.Join(wd, path)
	}
	if bytes, err := os.ReadFile(path); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}
