BIN := vct
VERSION := $(shell git describe --tags --always --dirty)

.PHONY: all $(BIN) test clean

all: $(BIN)

$(BIN):
	go build -ldflags='-s -w -X main.version=$(VERSION)' -trimpath

test:
	go test -v ./...

clean:
	rm -f $(BIN)
