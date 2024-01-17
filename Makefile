BIN := vct
VERSION := $(shell git describe --tags --always --dirty)

.PHONY: all $(BIN) clean

all: $(BIN)

$(BIN):
	go build -ldflags='-s -w -X main.version=$(VERSION)'

clean:
	rm -f $(BIN)
