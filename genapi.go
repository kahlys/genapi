package genapi

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
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

// Generate ...
func (d RestAPI) Generate(dir string) error {
	pkgname := strings.ToLower(d.ServiceName)

	// create destination directory
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		return fmt.Errorf("directory %v already exists", dir)
	}

	// Service
	service, err := writeService(pkgname, d.ServiceName, d.Endpoints...)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(fmt.Sprintf("%v/service.go", dir), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	fhandler, err := os.OpenFile(fmt.Sprintf("%v/handler.go", dir), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	t := Structure{
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
			Function: Function{
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
		endpoints = append(endpoints, fmt.Sprintf("r.HandleFunc(\"%s\", %s.%s).Methods(\"%v\")", m.URL, recvName, funcname, m.Method))
		meth := Method{
			Recv: fmt.Sprintf("%v *%v", recvName, name),
			Function: Function{
				Comments: fmt.Sprintf("%v ...", funcname),
				Name:     funcname,
				Params: []Parameter{
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
	endpoints = append(endpoints, "return r")

	// write Handler method
	meth := Method{
		Recv: fmt.Sprintf("%v *%v", strings.Split(strings.ToLower(name), "")[0], name),
		Function: Function{
			Comments: fmt.Sprintf("Handler returns the %v HTTP Handler.", name),
			Name:     "Handler",
			Content:  endpoints,
			Res:      []Parameter{{Type: "http.Handler"}},
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

// Method represents a method signature.
type Method struct {
	Recv string
	Function
}

// Function represents a function signature.
type Function struct {
	Comments string
	Name     string
	Params   []Parameter
	Res      []Parameter
	Content  []string
}

// Parameter represents a parameter in a function or method signature.
type Parameter struct {
	Name string
	Type string
}

// Structure represents a type
type Structure struct {
	Comments string
	Name     string
	Fields   []Parameter
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
