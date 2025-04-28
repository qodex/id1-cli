# id1-cli
Command Line Interface for id1 API

Build:

    go build -o id1

Usage: 

    id1 <options> [command] (data...)

Options:

    - env [set|get|del]: edit .env file
    - url: API endpoint URL
    - id: id1 id
    - create: create new id
    - key: path to private key pem file
    - enc: apply payload encoding
    - connect: connect to websocket
    - sync: apply key changes to workdir while connected

### Examples

Setup .env, create id:

    id1 env set url=http://localhost:8080
    id1 create test1 > test1.pem
    id1 env set id=test1
    id1 env set key=test1.pem

Read, Write keys:

    id1 set:/test1/pub/name Test One
    id1 set:/test1/pub/image < avatar.jpg
    id1 get:/test1/pub/name
    id1 get:/test1/pub/image > image.jpg
    id1 del:/test1/pub/image

Connect and sync:

    id1 connect sync

Set and delete after 60 seconds:

    id1 set:/test1/temp?ttl=60 Delete Me

Schedule a future command:

    id1 set:/test1/.after.1745816286236 set:/test1/alert1 Future is now
    
Schedule archive:

    id1 set:/test1/.after.1745816286236 mov:/test1/alert1 test1/archive/alert1 

List first 100 keys with values under 1KB:

    id1 get:/test1/*?recursive=true&size-limit=1024&limit=100&keys=true
    
    
