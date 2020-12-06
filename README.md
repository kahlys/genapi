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

Create a new project

```bash
genapi create <project_name>
cd <project_name>
```

It will create a directory with a configuration file you will need to write.

```yml
ServiceName: "<project_name>"
Endpoints:
  - Name: "GetElem"
    URL: "/api/elem/{id}"
    Method: "GET"
  - Name: "SetElem"
    URL: "/api/elem/{id}"
    Method: "POST"
```

In the project directory, after you write the configuration file, generate project files.

```bash
genapi init
```
