default: build

fmt:
	gofmt -w .

build:
	go build -o clilogin main.go
