.PHONY: all test
all:
	make -C cmd/gofrd
test:
	make -C lib/gofr test
	make -C cmd/gofrd test
