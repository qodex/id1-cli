package main

var man = `Usage: id1 <options> [command] (data...)

Options:
	- env [get|set|del]: .env config
	- dir: work dir
	- url: API endpoint URL
	- id: id1 id
	- create: id to create
	- key: path to a key file (private if connect, public if create)
	- mon: websocket to stdin, stdin to websocket
	- apply: apply stdin to work dir
	- watch: watch workdir, send events to stdout

Example:

id1 env set url=http://localhost:8080
id1 create test1 > test.pem
id1 env set id=test1
id1 env set key=test.pem

id1 set:/test1/pub/name Test One
id1 get:/test1/pub/name

id1 set:/test1/pub/image < avatar.jpg
id1 get:/test1/pub/image > image.jpg

Watch work dir, send changes to server:

id1 watch | id1 mon

Monitor server, apply events to work dir:

id1 mon | id1 apply
`
