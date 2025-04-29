package main

import (
	"bufio"
	"fmt"
	"os"
)

func scanData() []byte {
	if hasStdin() {
		scanner := bufio.NewScanner(os.Stdin)
		data := []byte{}
		for scanner.Scan() {
			data = append(data, scanner.Bytes()...)
			data = append(data, '\n')
		}
		if scanner.Err() != nil {
			fmt.Println(scanner.Err())
		}
		return data
	} else {
		return []byte{}
	}
}

func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if stat.Mode()&os.ModeCharDevice == 0 && stat.Size() > 0 {
		return true
	}
	return false
}
