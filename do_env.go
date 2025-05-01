package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/qodex/ff"
)

func env(args id1Args) {
	if f, ok := envOpFunc[args.envOp]; ok {
		f(args.KeyVal(args.envOp, "", ""))
	}
}

var envOpFunc = map[string]func(key, value string){
	"set": func(key, val string) {
		if len(key) == 0 {
			fmt.Println(man)
			os.Exit(0)
		}
		varName := strings.ToUpper(fmt.Sprintf("ID1_%s", key))
		ff.NewFsProps(".env").Set(varName, val)
	},

	"get": func(key, val string) {
		fmt.Println(os.Getenv(strings.ToUpper(fmt.Sprintf("ID1_%s", key))))
	},

	"del": func(key, val string) {
		if len(key) == 0 {
			fmt.Println(man)
		}
		varName := strings.ToUpper(fmt.Sprintf("ID1_%s", key))
		ff.NewFsProps(".env").Delete(varName)
	},

	"": func(key, val string) {
		for _, line := range os.Environ() {
			if strings.HasPrefix(line, "ID1_") {
				fmt.Println(strings.ToLower(strings.Replace(line, "ID1_", "", 1)))
			}
		}
	},
}
