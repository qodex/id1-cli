package main

import (
	"bufio"
	"os"

	id1 "github.com/qodex/id1-client-go"
)

func scanCommands(cmdOut chan id1.Command) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0), 10*id1.MB)

	data := []byte{}
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
		data = append(data, byte('\n'))
		if len(data) > 2 && string(data[len(data)-2:]) == "\n\n" {
			if cmd, err := id1.ParseCommand(data[:len(data)-2]); err == nil {
				cmdOut <- cmd
			}
			data = []byte{}
		}
	}
}
