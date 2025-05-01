# id1-cli
Command Line Interface for id1 API

Build:

    go build -o id1

Usage: 

    id1 <options> [command] (data...)

Options:

    - dir: local work dir
    - url: API endpoint URL
    - id: id1 id
    - key: path to key file
    - enc: apply payload encoding

Commands:

    - env: [get|set|del] .env options
    - create: create new id
    - serve: serve work dir
    - mon: monitor server events
    - watch: watch work dir events
    - apply: apply stdin commands to work dir

Key Commands:

    <operation>:<key>?[args] [data]

### Examples

Serve work dir on port:

    id1 env set dir=id1db
    id1 serve 8080

Setup .env, create id:

    id1 env set url=http://localhost:8080
    id1 create test1 > test1.pem
    id1 env set id=test1
    id1 env set key=test1.pem

Sync work dir changes to server:

    id1 watch | id1 mon

Apply server events to work dir:

    id1 mon | id1 apply

Read, Write keys:

    id1 set:/test1/pub/name Test One
    id1 set:/test1/pub/image < avatar.jpg
    id1 get:/test1/pub/name
    id1 get:/test1/pub/image > image.jpg
    id1 del:/test1/pub/image

Set and delete after 60 seconds:

    id1 set:/test1/temp\?ttl=60 Delete Me

Schedule future command:

    echo "set:/test1/alert1\nFuture is now" | id1 set:/test1/.after.1745816286237
    
Schedule archive:

    echo "mov:/test1/alert1\ntest1/arch/alert1" | id1 set:/test1/.after.1745816286237

List first 100 keys with values under 1KB:

    id1 get:/test1/\*\?recursive=true\&size-limit=1024\&limit=100\&keys=true
    
    
