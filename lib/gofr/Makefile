.PHONY: all bench clean lib

all: clean lib

clean:
	go clean .
	rm -f *.out

lib:
	go build .

test: clean
	go test -v -memprofile mem.out -cpuprofile cpu.out .

bench: lib
	go test -bench -v -memprofile mem.out -cpuprofile cpu.out .

graphs: bench
	go tool pprof --png gofr.test cpu.out > cpu_graph.png
	go tool pprof --png gofr.test mem.out > mem_graph.png

