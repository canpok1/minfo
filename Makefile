BINARY_NAME=minfo

run:
	@go run main.go ${server} ${option}

build:
	go build -o ${BINARY_NAME} main.go

build-release:
	goreleaser build --snapshot --clean

clean:
	go clean
	rm -f ${BINARY_NAME}
	rm -rf ./dist

test:
	go test -v ./...
