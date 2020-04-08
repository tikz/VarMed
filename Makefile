all: build

build: dep
	echo "Building VarQ web frontend"
	yarn --cwd web/ install
	yarn --cwd web/ build
	echo "Building VarQ to bin/varq"
	go build -o bin/varq -v

test: dep
	go test ./... -v

lint: dep
	golint

compile: dep
	echo "Cross compiling for all OSes and platforms"
	GOOS=linux GOARCH=amd64 go build -o bin/varq-linux-64
	GOOS=darwin GOARCH=amd64 go build -o bin/varq-darwin-64
	GOOS=windows GOARCH=amd64 go build -o bin/varq-windows-64.exe

dep:
	go get -u golang.org/x/lint/golint
	go get -v -d ./...
