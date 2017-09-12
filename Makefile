DEPS := github.com/google/uuid
DEPS += github.com/nfnt/resize
DEPS += github.com/lucasb-eyer/go-colorful

.PHONY: all test

all: commands test

clean:
	rm -fr vendor
	make -C cmd/gofrd clean
	make -C lib/gofr clean

vendor:
	mkdir -p vendor
	for repo in $(DEPS); do git clone https://$$repo vendor/$$repo; done
	rm -fr vendor/*/*/*/.git
	
commands:
	make -C cmd/gofrd

test:
	make -C lib/gofr test
	make -C cmd/gofrd test

