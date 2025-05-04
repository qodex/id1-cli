module github.com/qodex/id1-cli

go 1.23.3

// replace github.com/qodex/id1 => ../id1

// replace github.com/qodex/id1-client-go => ../id1-client-go

// replace github.com/qodex/ff => ../ff

require (
	github.com/qodex/ff v1.0.1
	github.com/qodex/id1-client-go v1.0.1
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	golang.org/x/sys v0.32.0 // indirect
)

require (
	github.com/fsnotify/fsnotify v1.9.0
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/joho/godotenv v1.5.1
	github.com/qodex/id1 v1.0.0
)
