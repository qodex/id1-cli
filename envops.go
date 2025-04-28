package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/qodex/ff"
)

var envOpFunc = map[string]func(key, value string){
	"set": func(key, val string) {
		if len(key) == 0 {
			fmt.Println(man)
			os.Exit(0)
		}
		varName := strings.ToUpper(fmt.Sprintf("ID1_%s", key))
		ff.NewFsProps(".env").Set(varName, val)
		os.Exit(0)
	},

	"get": func(key, val string) {
		fmt.Println(os.Getenv(strings.ToUpper(fmt.Sprintf("ID1_%s", key))))
		os.Exit(0)
	},

	"del": func(key, val string) {
		if len(key) == 0 {
			fmt.Println(man)
			os.Exit(0)
		}
		varName := strings.ToUpper(fmt.Sprintf("ID1_%s", key))
		ff.NewFsProps(".env").Delete(varName)
		os.Exit(0)
	},

	"": func(key, val string) {
		for _, line := range os.Environ() {
			if strings.HasPrefix(line, "ID1_") {
				fmt.Println(strings.ToLower(strings.Replace(line, "ID1_", "", 1)))
			}
		}
		os.Exit(0)
	},
}
