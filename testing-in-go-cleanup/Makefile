.PHONY: build
build:
		go build -o bin/server .

.PHONY: run
run: build
		./bin/server

.PHONY: test
test:
		go test ./... -v -count=1
