# To build the server
build:
	go build -o paseto-play cmd/paseto-play/main.go

# To generate a 64byte hash of your secret
keygen:
	go run cmd/keygen/main.go
