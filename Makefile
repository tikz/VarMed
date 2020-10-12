all: build

build: dep
	echo "Building VarMed web frontend"
	yarn --cwd web/ install
	yarn --cwd web/ build
	echo "Building VarMed binary to ./varmed"
	go build -o varmed -v

test: dep
	go test ./... -v

lint: dep
	golint

compile: dep
	echo "Cross compiling for all OSes and platforms"
	GOOS=linux GOARCH=amd64 go build -o dist/varmed-linux-64
	GOOS=darwin GOARCH=amd64 go build -o dist/varmed-darwin-64
	GOOS=windows GOARCH=amd64 go build -o dist/varmed-windows-64.exe

dep:
	go get -u golang.org/x/lint/golint
	go get -v -d ./...
