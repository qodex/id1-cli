package main

var man = `Usage: id1 <options> [command] (data...)

Options:
	- url: API endpoint URL
	- id: id1 id
	- create: id to create
	- key: path to a key file (private if connect, public if create)
	- sync: apply incoming commands to workdir

Environment:
	ID1_URL: Default id1 API endpoint url
	ID1_ID: Default id1 id
	ID1_KEY_PATH: Default key path

Example:

id1 env set url=http://localhost:8080
id1 create test1 > test.pem
id1 env set id=test1
id1 env set key=test.pem

id1 set:/test1/pub/name Test One
id1 get:/test1/pub/name

id1 set:/test1/pub/image < avatar.jpg
id1 get:/test1/pub/image > image.jpg
`
