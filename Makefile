BIN   := $(shell basename $(CURDIR))
LIBS  := $(shell find lib -type d -mindepth 1 -maxdepth 1)
CMDS  := $(shell find cmd -type d -mindepth 1 -maxdepth 1)

TOOLS := github.com/golang/dep/cmd/dep

.PHONY: all clean cmd docker run test tools

all: test

clean:
	for d in $(CMDS) $(LIBS) ; do make -C $$d clean ; done
	rm -fr vendor

tools:
	for t in $(TOOLS); do go get $$t; done

vendor: tools
	dep ensure

cmd: vendor
	for d in $(CMDS) ; do make -C $$d $$d; done

docker: vendor
	make -C cmd/gofrd docker
	docker-compose up --build -d

test: vendor
	for d in $(LIBS) $(CMDS); do make -C $$d test ; done

run: vendor
	make -C cmd/gofrd run

