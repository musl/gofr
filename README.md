# gofr

This is a rendering service and UI for browsing the fractal-y goodnees of the
Mandelbrot set.

## Build & Play

1. Have a go directory and properly set GOPATH. See: [https://golang.org/doc/code.html](https://golang.org/doc/code.html)
1. `go get github.com/musl/gofr`
1. `cd $GOPATH/src/github.com/musl/gofr`
1. `make`
1. Browse to: [http://127.0.0.1:8000/](http://127.0.0.1:8000/)

## Vendored Libraries

- Go
    - [https://github.com/nfnt/resize](https://github.com/nfnt/resize)
    - [https://github.com/google/uuid](https://github.com/google/uuid)
- JS
    - [http://www.ractivejs.org/](http://www.ractivejs.org/)
    - [https://ace.c9.io/](https://ace.c9.io/)
    - [https://jquery.ycom/](https://jquery.ycom/)
- CSS
    - [http://purecss.io/](http://purecss.io/)
    - [http://fontawesome.io/](http://fontawesome.io/)

## Updating Vendored Go Dependencies

1. `make clean vendor`
1. `git add -f vendor`
1. `git ci -m "Update vendored dependencies"`

