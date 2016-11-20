.PHONY: all bench clean lib

all: clean lib

clean:
	go clean .
	rm -f $(IMAGE_FILE)
	rm -f *.out

lib:
	go build .

test: clean
	go test -memprofile mem.out -cpuprofile cpu.out .

graphs: test
	go tool pprof --png $(BIN).test cpu.out > cpu_graph.png
	go tool pprof --png $(BIN).test mem.out > mem_graph.png

bench: $(BIN)
	go test -bench .
