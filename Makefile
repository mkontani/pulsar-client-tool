BINARY := pulsar-client-tool
VERSION ?= dev
LDFLAGS := -ldflags "-X github.com/mkontani/pulsar-client-tool/cmd.Version=$(VERSION)"

.PHONY: build test lint clean vet fmt

build:
	go build $(LDFLAGS) -o $(BINARY) .

test:
	go test ./...

lint:
	golangci-lint run

vet:
	go vet ./...

fmt:
	gofmt -w .

clean:
	rm -f $(BINARY)
