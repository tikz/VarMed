all: build

build: dep
	echo "Building RespDB web frontend"
	yarn --cwd web/ install
	yarn --cwd web/ build
	echo "Building RespDB binary to ./respdb"
	go build -o respdb -v

test: dep
	go test ./... -v

lint: dep
	golint

compile: dep
	echo "Cross compiling for all OSes and platforms"
	GOOS=linux GOARCH=amd64 go build -o dist/respdb-linux-64
	GOOS=darwin GOARCH=amd64 go build -o dist/respdb-darwin-64
	GOOS=windows GOARCH=amd64 go build -o dist/respdb-windows-64.exe

dep:
	go get -u golang.org/x/lint/golint
	go get -v -d ./...
