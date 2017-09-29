BIN := $(shell basename $(CURDIR))

.PHONY: all clean clobber run test

all: clean frontend run

clean:
	rm -f $(BIN)
	rm -fr build

clobber: clean
	rm -fr node_modules

node_modules:
	npm i

frontend: node_modules
	webpack

$(BIN): 
	go build .

test: $(BIN)
	go test -v .

run: $(BIN)
	./$(BIN)

