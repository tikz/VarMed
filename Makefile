dep:
	go get -v -d ./...

build:
	echo "Building VarQ to bin/varq"
	go build -o bin/varq -v

test:
	go test ./... -v

lint:
	golint

compile:
	echo "Cross compiling for all OSes and platforms"
	GOOS=linux GOARCH=amd64 go build -o bin/varq-linux-64
	GOOS=darwin GOARCH=amd64 go build -o bin/varq-darwin-64
	GOOS=windows GOARCH=amd64 go build -o bin/varq-windows-64.exe

all: build