package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/qodex/id1"
)

func serve(path, port string) {
	ctx := context.Background()
	http.HandleFunc("/{key...}", id1.Handle(path, ctx))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("error starting service: %s", err)
	}
	ctx.Done()
}
