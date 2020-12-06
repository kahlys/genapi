package genapi

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// RestAPI description
type RestAPI struct {
	ServiceName string
	Endpoints   []Endpoint
}

// Endpoint http rest api
type Endpoint struct {
	Name   string
	URL    string
	Method string
}

// Generate the output code in given directory
func (d RestAPI) Generate() error {
	pkgname := strings.ToLower(d.ServiceName)

	// Service
	service, err := writeService(pkgname, d.ServiceName, d.Endpoints...)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(".", "service.go"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.Write(service)
	if err != nil {
		return err
	}

	// Handler
	handler, err := writeHandler(pkgname, d.ServiceName, d.Endpoints...)
	if err != nil {
		return err
	}
	fhandler, err := os.OpenFile(filepath.Join(".", "handler.go"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	t := structure{
		Comments: fmt.Sprintf("%v ...", name),
		Name:     name,
	}
	if err := tmplTypes.Execute(buf, t); err != nil {
		return []byte{}, err
	}

	// write methods
	for _, m := range methods {
		funcname := strings.Title(strings.ToLower(m.Name))
		meth := method{
			Recv: fmt.Sprintf("%v *%v", strings.Split(strings.ToLower(name), "")[0], name),
			function: function{
				Comments: fmt.Sprintf("%v ...", funcname),
				Name:     funcname,
				Res:      []parameter{{"err", "error"}},
				Content:  []string{"return nil"},
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
		endpoints = append(endpoints, fmt.Sprintf("r.HandleFunc(\"%s\", %s.%s).Methods(\"%v\")", m.URL, recvName, funcname, m.Method))
		meth := method{
			Recv: fmt.Sprintf("%v *%v", recvName, name),
			function: function{
				Comments: fmt.Sprintf("%v ...", funcname),
				Name:     funcname,
				Params: []parameter{
					{Name: "w", Type: "http.ResponseWriter"},
					{Name: "req", Type: "*http.Request"},
				},
				Content: []string{
					fmt.Sprintf("err := %v.%v()", recvName, strings.Title(strings.ToLower(m.Name))),
					"if err != nil {http.Error(w, err.Error(),http.StatusInternalServerError,); return}",
					"w.WriteHeader(http.StatusNotImplemented)",
				},
			},
		}
		if err := tmplMethodes.Execute(buf, meth); err != nil {
			return []byte{}, err
		}
	}
	endpoints = append(endpoints, "return r")

	// write Handler method
	meth := method{
		Recv: fmt.Sprintf("%v *%v", strings.Split(strings.ToLower(name), "")[0], name),
		function: function{
			Comments: fmt.Sprintf("Handler returns the %v HTTP Handler.", name),
			Name:     "Handler",
			Content:  endpoints,
			Res:      []parameter{{Type: "http.Handler"}},
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

// method represents a method signature.
type method struct {
	Recv string
	function
}

// function represents a function signature.
type function struct {
	Comments string
	Name     string
	Params   []parameter
	Res      []parameter
	Content  []string
}

// Parameter represents a parameter in a function or method signature.
type parameter struct {
	Name string
	Type string
}

// Structure represents a type
type structure struct {
	Comments string
	Name     string
	Fields   []parameter
}

const textMethods = "{{if .Comments}}// {{.Comments}}\n{{end}}" +
	"func ({{.Recv}}) {{.Name}}" +
	"({{range .Params}}{{.Name}} {{.Type}}, {{end}})" +
	"({{range .Res}}{{.Name}} {{.Type}}, {{end}})" +
	"{\n {{if .Content}} {{range $i, $x := .Content}} {{$x}}; {{end}} {{else}}panic(\"not implemented\"){{end}}}\n\n"

const textTypes = "{{if .Comments}}// {{.Comments}}\n{{end}}" +
	"type {{.Name}} struct" +
	"{}\n\n"

var (
	tmplMethodes = template.Must(template.New("methods").Parse(textMethods))
	tmplTypes    = template.Must(template.New("types").Parse(textTypes))
)
