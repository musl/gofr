.PHONY: all clean clobber test

all: test

clean:
	make -C cmd/gofrd clean
	make -C lib/gofr clean
	rm -fr vendor
	
clobber: clean
	make -C cmd/gofrd clobber

run: vendor
	make -C cmd/gofrd run

docker: vendor
	make -C cmd/gofrd docker
	docker-compose up --build -d

test: vendor
	make -C lib/gofr test
	make -C cmd/gofrd test

vendor:
	dep ensure

