build:
	go build -o bin/gorate ./cmd/server

run: build
	./bin/gorate

test:
	go test ./... -race

clean:
	rm -rf bin/
