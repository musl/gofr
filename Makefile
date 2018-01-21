.PHONY: all clean clobber test

all: commands test

clean:
	make -C cmd/gofrd clean
	make -C lib/gofr clean
	
clobber: clean
	rm -fr vendor
	make -C cmd/gofrd clobber

commands: vendor
	make -C cmd/gofrd

daemon: vendor
	make -C cmd/gofrd clean run

docker: vendor
	make -C cmd/gofrd docker
	docker-compose up --build -d

test: vendor
	make -C lib/gofr test
	make -C cmd/gofrd test

vendor:
	dep ensure

