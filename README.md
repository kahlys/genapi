# kahlys/genapi

[![godoc](https://godoc.org/github.com/kahlys/genapi?status.svg)](https://godoc.org/github.com/kahlys/genapi)
[![build](https://api.travis-ci.org/kahlys/genapi.svg?branch=master)](https://travis-ci.org/kahlys/genapi)
[![go report](https://goreportcard.com/badge/github.com/kahlys/genapi)](https://goreportcard.com/report/github.com/kahlys/genapi)

Quick golang code generation for rest api microservices using gorilla mux.

# Installation

With a correctly configured [Go toolchain](https://golang.org/doc/install):

```sh
$ git clone github.com/kahlys/genapi/
$ cd proxy
$ make install
```

# Usage

Global configuration. The _service name_ will be the name of the structure for the microservices.

```sh
>>> config
Configure service
Service name: Booker
Destination directory: booker
```

Add endpoints. The _name_ will be used to give a name to function.

```sh
>>> add
Add an endpoint
Name: GetBook
URL: /api/book/{id}
Method: GET
Adding book : GET /api/book/{id}

>>> add
Add an endpoint
Name: SetBook
URL: /api/book
Method: POST
Adding setBook : POST /api/book
```

Run to generate your files. An error message will be printed if athe operation failed.

```sh
>>> run
```

Generated files in this example are

```go
// file service.go
package main

// Booker ...
type Booker struct{}

// GetBook ...
func (b *Booker) GetBook() {
    panic("not implemented")
}

// SetBook ...
func (b *Booker) SetBook() {
    panic("not implemented")
}
```

```go
// file handler.go
package main

import (
    "net/http"

    "github.com/gorilla/mux"
)

// handleGetBook ...
func (b *Booker) handleGetBook(w http.ResponseWriter, req *http.Request) {
    b.GetBook()
}

// handleSetBook ...
func (b *Booker) handleSetBook(w http.ResponseWriter, req *http.Request) {
    b.SetBook()
}

// Handler returns the Booker HTTP Handler.
func (b *Booker) Handler() http.Handler {
    r := mux.NewRouter()
    r.HandleFunc("/api/book/{id}", b.handleGetBook).Methods("GET")
    r.HandleFunc("/api/book", b.handleSetBook).Methods("POST")
    return r
}
```
