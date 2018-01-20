BIN := $(shell basename $(CURDIR))

.PHONY: all clean clean_build clobber run test

all: clean clean_build frontend run

clean:
	rm -f $(BIN)

clean_build:
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

docker: clean frontend
	GOOS=linux GOARCH=amd64 go build .

