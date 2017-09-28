BIN := $(shell basename $(CURDIR))

.PHONY: all clean run

all: clean frontend run

clean:
	rm -f $(BIN)
	rm -fr build

dist_clean: clean
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

