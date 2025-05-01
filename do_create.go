package main

import (
	"fmt"
	"os"

	id1 "github.com/qodex/id1-client-go"
)

func createId(id string, c id1.Id1Client) {
	pubKey := id1.KK(id, "pub", "key")
	if _, err := c.Get(pubKey); err == nil {
		fmt.Println("id already exists")
		os.Exit(1)
	}
	publicKey := ""
	privateKey := ""
	if keyPair, err := genKey(1024); err != nil {
		fmt.Println("key gen error", err)
		os.Exit(1)
	} else {
		publicKey = keyPair.public
		privateKey = keyPair.private
	}
	if err := c.Set(pubKey, []byte(publicKey)); err != nil {
		fmt.Printf("set %s error %s\n", pubKey, err)
		os.Exit(1)
	} else {
		fmt.Print(privateKey)
		os.Exit(0)
	}
}
