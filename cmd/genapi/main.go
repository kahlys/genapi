package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"os"
	"strings"

	"github.com/abiosoft/ishell"
)

var version = "undefined"

var datas struct {
	directory   string
	serviceName string
	endpoints   []Endpoint
}

// TODO add database structure (optional)
// TODO add main file
// TODO structure with dir cmd and dir pkg

// TODO make configurable
const (
	pkgname = "main"
)

func main() {
	shell := ishell.New()

	// display info.
	shell.Println("Genapi - version", version)

	// configure
	shell.AddCmd(&ishell.Cmd{
		Name: "config",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)
			c.Println("Configure service")
			c.Print("Service name: ")
			servicename := c.ReadLine()
			c.Print("Destination directory: ")
			dir := c.ReadLine()
			// TODO verify name (no space, etc)
			datas.serviceName = servicename
			datas.directory = dir
		},
		Help: "global configuration",
	})

	// add endpoint
	shell.AddCmd(&ishell.Cmd{
		Name: "add",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)
			c.Println("Add an endpoint")
			e := Endpoint{}
			c.Print("Name: ")
			e.Name = c.ReadLine()
			c.Print("URL: ")
			e.URL = c.ReadLine()
			c.Print("Method: ")
			e.Method = c.ReadLine()
			// TODO verify endpoint before adding it to the list
			c.Println("Adding", e.Name, ":", e.Method, e.URL)
			datas.endpoints = append(datas.endpoints, e)
		},
		Help: "add an endpoint",
	})

	// show configuration
	shell.AddCmd(&ishell.Cmd{
		Name: "show",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)
			// TODO proper display
			c.Println(datas)
		},
		Help: "show configuration",
	})

	// show configuration
	shell.AddCmd(&ishell.Cmd{
		Name: "run",
		Func: func(c *ishell.Context) {
			c.ShowPrompt(false)
			defer c.ShowPrompt(true)
			if err := generate(); err != nil {
				c.Println("ERROR:", err)
			}
		},
		Help: "generate golang code",
	})

	shell.Run()
	shell.Close()
}

func generate() error {
	// create destination directory
	if err := os.Mkdir(datas.directory, os.ModePerm); err != nil {
		return fmt.Errorf("directory %v already exists", datas.directory)
	}

	// Service
	service, err := writeService(pkgname, datas.serviceName, datas.endpoints...)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(fmt.Sprintf("%v/service.go", datas.directory), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.Write(service)
	if err != nil {
		return err
	}

	// Handler
	handler, err := writeHandler(pkgname, datas.serviceName, datas.endpoints...)
	if err != nil {
		return err
	}
	fhandler, err := os.OpenFile(fmt.Sprintf("%v/handler.go", datas.directory), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer fhandler.Close()
	if err != nil {
		return err
	}
	_, err = fhandler.Write(handler)
	if err != nil {
		return err
	}

	return nil
}

func writeService(pkg, name string, methods ...Endpoint) ([]byte, error) {
	pkg = strings.ToLower(pkg)
	name = strings.Title(strings.ToLower(name))

	// package header
	buf := bytes.NewBufferString(fmt.Sprintf("package %v\n", pkg))

	// write structure
	t := Type{
		Comments: fmt.Sprintf("%v ...", name),
		Name:     name,
	}
	if err := tmplTypes.Execute(buf, t); err != nil {
		return []byte{}, err
	}

	// write methods
	for _, m := range methods {
		funcname := strings.Title(strings.ToLower(m.Name))
		meth := Method{
			Recv: fmt.Sprintf("%v *%v", strings.Split(strings.ToLower(name), "")[0], name),
			Func: Func{
				Comments: fmt.Sprintf("%v ...", funcname),
				Name:     funcname,
			},
		}
		if err := tmplMethodes.Execute(buf, meth); err != nil {
			return []byte{}, err
		}
	}

	// format
	buf = bytes.NewBufferString(strings.Replace(buf.String(), "&#34;", "\"", -1)) // TODO find a way to use quotes in template
	pretty, err := format.Source(buf.Bytes())
	if err != nil {
		return []byte{}, err
	}

	return pretty, nil
}

func writeHandler(pkg, name string, methods ...Endpoint) ([]byte, error) {
	pkg = strings.ToLower(pkg)
	name = strings.Title(strings.ToLower(name))

	// package header
	buf := bytes.NewBufferString(fmt.Sprintf("package %v\n", pkg))
	buf.WriteString("import( \"net/http\" \n\n \"github.com/gorilla/mux\" )\n")

	// write methods
	endpoints := []string{"r := mux.NewRouter()"}
	for _, m := range methods {
		funcname := fmt.Sprintf("handle%v", strings.Title(strings.ToLower(m.Name)))
		recvName := strings.Split(strings.ToLower(name), "")[0]
		endpoints = append(endpoints, fmt.Sprintf("r.HandleFunc(\"%s\", %s.%s)", m.URL, recvName, funcname))
		meth := Method{
			Recv: fmt.Sprintf("%v *%v", recvName, name),
			Func: Func{
				Comments: fmt.Sprintf("%v ...", funcname),
				Name:     funcname,
				Params: []Param{
					{Name: "w", Type: "http.ResponseWriter"},
					{Name: "req", Type: "*http.Request"},
				},
				Content: []string{fmt.Sprintf("%v.%v()", recvName, strings.Title(strings.ToLower(m.Name)))},
			},
		}
		if err := tmplMethodes.Execute(buf, meth); err != nil {
			return []byte{}, err
		}
	}
	endpoints = append(endpoints, fmt.Sprintf("return r"))

	// write Handler method
	meth := Method{
		Recv: fmt.Sprintf("%v *%v", strings.Split(strings.ToLower(name), "")[0], name),
		Func: Func{
			Comments: fmt.Sprintf("Handler returns the %v HTTP Handler.", name),
			Name:     "Handler",
			Content:  endpoints,
			Res:      []Param{{Type: "http.Handler"}},
		},
	}
	if err := tmplMethodes.Execute(buf, meth); err != nil {
		return []byte{}, err
	}

	// format
	buf = bytes.NewBufferString(strings.Replace(buf.String(), "&#34;", "\"", -1)) // TODO find a way to use quotes in template
	pretty, err := format.Source(buf.Bytes())
	if err != nil {
		return []byte{}, err
	}

	return pretty, nil
}

// Endpoint http rest api
type Endpoint struct {
	Name   string
	URL    string
	Method string
}

// Method represents a method signature.
type Method struct {
	Recv string
	Func
}

// Func represents a function signature.
type Func struct {
	Comments string
	Name     string
	Params   []Param
	Res      []Param
	Content  []string
}

// Param represents a parameter in a function or method signature.
type Param struct {
	Name string
	Type string
}

// Type represents a type
type Type struct {
	Comments string
	Name     string
	Fields   []Param
}

const textMethods = "{{if .Comments}}// {{.Comments}}\n{{end}}" +
	"func ({{.Recv}}) {{.Name}}" +
	"({{range .Params}}{{.Name}} {{.Type}}, {{end}})" +
	"({{range .Res}}{{.Name}} {{.Type}}, {{end}})" +
	"{\n {{if .Content}} {{range $i, $x := .Content}} {{$x}}; {{end}} {{else}}panic(\"not implemented\"){{end}}}\n\n"

var tmplMethodes = template.Must(template.New("methods").Parse(textMethods))

const textTypes = "{{if .Comments}}// {{.Comments}}\n{{end}}" +
	"type {{.Name}} struct" +
	"{}\n\n"

var tmplTypes = template.Must(template.New("types").Parse(textTypes))
