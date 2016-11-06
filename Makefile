.PHONY: all bench clean

all: clean test

clean:
	go clean .
	rm -f $(IMAGE_FILE)
	rm -f *.out

$(BIN):
	go build .

test: clean
	go test -memprofile mem.out -cpuprofile cpu.out .

graphs: test
	go tool pprof --png $(BIN).test cpu.out > cpu_graph.png
	go tool pprof --png $(BIN).test mem.out > mem_graph.png

bench: $(BIN)
	go test -bench .
