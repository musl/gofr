DEPS := github.com/google/uuid
DEPS += github.com/nfnt/resize
DEPS += github.com/lucasb-eyer/go-colorful
DEPS += github.com/stretchr/testify

.PHONY: all clean clobber test vendor

all: commands test

clean:
	rm -fr vendor
	make -C cmd/gofrd clean
	make -C lib/gofr clean
	
clobber: clean
	make -C cmd/gofrd clobber

vendor:
	mkdir -p vendor
	for repo in $(DEPS); do git clone https://$$repo vendor/$$repo; done
	rm -fr vendor/*/*/*/.git
	
commands:
	make -C cmd/gofrd

test:
	make -C lib/gofr test
	make -C cmd/gofrd test

daemon:
	make -C cmd/gofrd clean run

docker:
	make -C cmd/gofrd docker
	docker-compose up --build -d

