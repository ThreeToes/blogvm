build: test
	go build -o ./bin/blogvm ./cmd/main.go

test:
	go test -failfast ./...

clean:
	rm -rf ./bin