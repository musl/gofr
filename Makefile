.PHONY: all clean clobber test

all: vendor commands test

clean:
	make -C cmd/gofrd clean
	make -C lib/gofr clean
	
clobber: clean
	rm -fr vendor
	make -C cmd/gofrd clobber

commands:
	make -C cmd/gofrd

daemon:
	make -C cmd/gofrd clean run

docker:
	make -C cmd/gofrd docker
	docker-compose up --build -d

test:
	make -C lib/gofr test
	make -C cmd/gofrd test

vendor:
	dep ensure

