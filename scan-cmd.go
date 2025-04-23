package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	id1 "github.com/qodex/id1-client-go"
)

func scanCommands(cmdCh chan id1.Command) {
	scanner := bufio.NewScanner(os.Stdin)
	str := ""
	for scanner.Scan() {
		str += scanner.Text() + "\n"
		if strings.HasSuffix(str, "\n\n") {
			str = strings.TrimSuffix(str, "\n\n")
			str = strings.TrimPrefix(str, "\n")
			str = strings.TrimPrefix(str, "\n")
			str = strings.TrimPrefix(str, "\n")
			if cmd, err := id1.ParseCommand([]byte(str)); err == nil {
				cmdCh <- cmd
			}
			str = ""
		}
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
		os.Exit(0)
	}
}
