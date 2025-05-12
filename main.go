package main

import "fmt"

func main() {
	args := getArgs()
	switch {
	case args.env:
		env(args)
	case args.create:
		if c, err := getClient(args); err == nil {
			createId(args.createId, *c)
		}
	case args.serve:
		serve(args.dir, args.port)
	case args.apply:
		apply(args)
	case args.watch:
		watch(args)
	case args.mon:
		if c, err := getClient(args); err == nil {
			mon(args, *c)
		}
	case args.filter:
		filter(args)
	case args.cmd != nil:
		if c, err := getClient(args); err == nil {
			cmdExec(*args.cmd, *c)
		}
	default:
		fmt.Println(man)
	}
}
