run:
	@go run main.go ${server} ${option}

build:
	goreleaser build --snapshot --clean

clean:
	rm -rf ./dist

test:
	go test -v ./...
