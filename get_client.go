package main

import (
	"encoding/base64"
	"fmt"

	id1 "github.com/qodex/id1-client-go"
)

func getClient(args id1Args) (*id1.Id1Client, error) {
	c, err := id1.NewHttpClient(args.url)
	if err != nil {
		return nil, err
	}

	switch args.enc {
	case "base64":
		prep := []func(cmd *id1.Command) error{
			func(cmd *id1.Command) error {
				if len(cmd.Data) > 0 {
					cmd.Data = []byte(base64.StdEncoding.EncodeToString(cmd.Data))
				}
				return nil
			},
		}
		postp := []func(data []byte, err error) ([]byte, error){}
		proxy := id1.NewId1ClientProxy(c, prep, postp)
		c = proxy
	}

	if len(args.id) > 0 && len(args.key) > 0 {
		if err := c.Authenticate(args.id, args.key); err != nil {
			fmt.Printf("authentication failed: %s %s %d %s\n", err, args.id, len(args.key), args.url)
			return nil, fmt.Errorf("auth failed")
		}
	}

	return &c, nil
}
