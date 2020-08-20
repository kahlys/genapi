# kahlys/genapi

[![PkgGoDev](https://pkg.go.dev/badge/github.com/kahlys/genapi)](https://pkg.go.dev/github.com/kahlys/genapi)
![build](https://github.com/kahlys/genapi/workflows/build/badge.svg)
[![go report](https://goreportcard.com/badge/github.com/kahlys/genapi)](https://goreportcard.com/report/github.com/kahlys/genapi)

Quick golang code generation for rest api microservices using gorilla mux.

## Installation

With a correctly configured [Go toolchain](https://golang.org/doc/install):

```bash
go get -u github.com/kahlys/genapi
```

## Usage

Write Rest API description in a configuration file.

```yml
ServiceName: "booker"
Endpoints:
  - Name: "GetBook"
    URL: "/api/book/{id}"
    Method: "GET"
  - Name: "SetBook"
    URL: "/api/book"
    Method: "POST"
```

Run to generate your files by giving the path to the configuration file, and the output directory path.

```bash
genapi -config ./example/config.yml - dir mydir
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
