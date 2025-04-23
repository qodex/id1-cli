# id1-cli
Command Line Interface for id1 API

go build -o id1


    Usage: 

    id1 <options> [command] (data...)

    Options:

        -url: API endpoint URL
        -id: id1 id
        -create: id to create
        -key: path to a key file (private if connect, public if create)

    Environment:
    
    ID1_URL: Default id1 API endpoint url
    ID1_ID: Default id1 id
    ID1_KEY_PATH: Default key path

    Example:

    id1 -url http://localhost:8080 get:/test1/pub/key
    id1 -url http://localhost:8080 -create test1 > test1.pem

    export ID1_URL=http://localhost:8080
    export ID1_ID=test1
    export ID1_KEY_PATH=test1.pem

    id1 set:/test1/pub/name Test One
    id1 get:/test1/pub/name

    id1 set:/test1/pub/image < avatar.jpg
    id1 get:/test1/pub/image > image.jpg