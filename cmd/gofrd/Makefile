BIN := $(shell basename $(CURDIR))

.PHONY: all clean clean_build clobber run test

all: test

clean:
	rm -f $(BIN)

clean_build:
	rm -fr build

clobber: clean clean_build
	rm -fr node_modules

node_modules:
	npm i

frontend: node_modules
	webpack

$(BIN): 
	go build .

test: $(BIN)
	go test -v .

run: $(BIN) clean_build frontend
	./$(BIN)

docker: clean frontend
	GOOS=linux GOARCH=amd64 go build .

