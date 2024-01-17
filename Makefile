BIN := vct

.PHONY: all $(BIN) clean

all: $(BIN)

$(BIN):
	go build -ldflags='-s -w'

clean:
	rm -f $(BIN)
