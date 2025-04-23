package main

import (
	"fmt"

	id1 "github.com/qodex/id1-client-go"
)

func createId(id string, c id1.Id1Client) (string, error) {
	pubKey := id1.KK(id, "pub", "key")
	if _, err := c.Get(pubKey); err == nil {
		return "", fmt.Errorf("id already exists")
	} else if key, err := genKey(1024); err != nil {
		return "", err
	} else if err := c.Set(pubKey, []byte(key.public)); err != nil {
		return "", err
	} else {
		return key.private, nil
	}
}
