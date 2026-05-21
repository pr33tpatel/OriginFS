build:
	go build -o bin/OriginFS

run: build
	./bin/OriginFS

test:
	go test ./... -v
