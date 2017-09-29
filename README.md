[![Build Status](https://travis-ci.org/musl/gofr.svg?branch=master)](https://travis-ci.org/musl/gofr)
[![GoDoc](https://godoc.org/github.com/musl/gofr/lib/gofr?status.svg)](http://godoc.org/github.com/musl/gofr/lib/gofr)
[![GoDoc](https://godoc.org/github.com/musl/gofr/cmd/gofrd?status.svg)](http://godoc.org/github.com/musl/gofr/cmd/gofrd)

# gofr

This is a rendering service and UI for browsing the fractal goodnees of the
Mandelbrot set.

## Play

1. Install global prerequisites:

  1. npm (`brew install npm`)
  1. webpack (`npm install -g webpack`)

1. Fetch the project: `go get github.com/musl/gofr/...`

    If you run into problems, check or set up your go directory and `GOPATH` with the instructions here: [https://golang.org/doc/code.html](https://golang.org/doc/code.html)

1. Run the service:

  1. `cd $GOPATH/src/github.com/musl/gofr`
  1. `make`

1. Browse to: [http://127.0.0.1:8000/](http://127.0.0.1:8000/)

## Update Vendored Go Dependencies

1. `make clean vendor`
1. `git add -f vendor`
1. `git ci -m "Update vendored dependencies"`

## List of Vendored Libraries

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

## Docs

- [lib/gofr](http://godoc.org/github.com/musl/gofr/lib/gofr)

- [cmd/gofrd](http://godoc.org/github.com/musl/gofr/cmd/gofrd)
    
    The binary is more or less a [12-factor app](http://12factor.net)
    that accepts configuration via the environment:

    - `GOFR_STATIC_DIR`: The path to the static assets for gofrd.  Default: `./static`
    - `GOFR_BIND_ADDR`: The address and port to bind to. Default: `0.0.0.0:8000`

