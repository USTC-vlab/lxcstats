SRC := $(wildcard *.go go.mod)
BIN := lxcstats

.PHONY: all

all: $(BIN)

$(BIN): $(SRC)
	go build -ldflags='-s -w'
