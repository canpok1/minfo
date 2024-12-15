BINARY_NAME=minfo

run:
	@go run cmd/main.go ${origin} ${limit}

build:
	go build -o ${BINARY_NAME} cmd/main.go

clean:
	go clean
	rm -f ${BINARY_NAME}

test:
	go test -v ./...
