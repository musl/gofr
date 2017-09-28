BIN := $(shell basename $(CURDIR))

.PHONY: all clean run

all: clean run

clean:
	rm -f $(BIN)
	rm -fr build

frontend:
	webpack

$(BIN): 
	go build .

test: $(BIN)
	go test -v .

run: $(BIN)
	./$(BIN)

