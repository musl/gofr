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

dev: clean frontend
	npm start

$(BIN): 
	go build .

test: $(BIN)
	go test -v .

run: $(BIN)
	$(CURDIR)/$(BIN)

docker: clean frontend
	GOOS=linux GOARCH=amd64 go build .

