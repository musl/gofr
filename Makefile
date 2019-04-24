.PHONY: all clean clobber test

all: test

clean:
	make -C cmd/gofrd clean
	make -C lib/gofr clean
	rm -fr vendor
	
clobber: clean
	make -C cmd/gofrd clobber

run:
	make -C cmd/gofrd run

docker:
	make -C cmd/gofrd docker
	docker-compose up --build -d

test:
	make -C lib/gofr test
	make -C cmd/gofrd test
