package main

import (
	"fmt"
	"strings"

	id1 "github.com/qodex/id1-client-go"
)

func parseCommand(str string, data []byte) (id1.Command, error) {
	if strings.IndexByte(str, '\n') < 0 && strings.IndexByte(str, ' ') > 0 {
		str = strings.Replace(str, " ", "\n", 1)
	}
	if cmd, err := id1.ParseCommand([]byte(str)); err != nil {
		return id1.Command{}, fmt.Errorf("invalid command")
	} else {
		data := append(cmd.Data, data...)
		cmd.Data = data
		return cmd, nil
	}
}
