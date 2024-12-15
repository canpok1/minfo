BINARY_NAME=minfo

run:
	@go run cmd/main.go

build:
	go build -o ${BINARY_NAME} cmd/main.go

clean:
	go clean
	rm -f ${BINARY_NAME}
