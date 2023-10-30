BIN_NAME=verit

build:
	CGO_ENABLED=0 go build -o dist/$(BIN_NAME) -ldflags "-s -w" main.go


tiny:
	tinygo build -o dist/$(BIN_NAME) main.go
	strip dist/$(BIN_NAME)