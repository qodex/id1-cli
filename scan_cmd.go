package main

import (
	"bytes"
	"encoding/base64"

	"github.com/qodex/ff"
	id1 "github.com/qodex/id1-client-go"
)

func scanCommands(cmdOut chan id1.Command, eof chan bool) {

	dataIn := make(chan []byte, 8)
	scanEof := make(chan bool)

	go ff.ScanStdin([]byte("\n\n"), dataIn, scanEof)

	for {
		select {
		case <-scanEof:
			eof <- true
			return
		case data := <-dataIn:
			data = bytes.TrimLeft(data, "\n")
			if cmd, err := id1.ParseCommand(data); err == nil {
				if cmd.Args["enc"] == "base64" {
					if data, err := base64.StdEncoding.DecodeString(string(cmd.Data)); err == nil {
						cmd.Data = data
					}
				}
				cmdOut <- cmd
			}
		}
	}
}
