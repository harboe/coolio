package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	router "github.com/julienschmidt/httprouter"
)

var (
	port        string
	proxy       string
	templateDir string
)

type context struct{}

func main() {
	flag.StringVar(&port, "port", "localhost:8080", "port")
	flag.StringVar(&proxy, "proxy", "", "used with ngrok")
	flag.StringVar(&templateDir, "template", "templates", "location of templates files, remember to run gulp.")
	flag.Parse()

	ctx := &context{}

	r := router.New()
	// content
	r.GET("/v1/:view", ctx.contentHandler)
	r.GET("/v1/:view/:version", ctx.contentHandler)
	r.GET("/v1/:view/:version/:file", ctx.contentHandler)

	// preview && save handler
	r.POST("/:view", ctx.saveHandler)

	// index
	r.GET("/", ctx.indexHandler)

	// setup static content.
	r.ServeFiles("/static/*filepath", http.Dir(templateDir+"/static"))
	r.ServeFiles("/fonts/*filepath", http.Dir(templateDir+"/static/fonts"))

	fmt.Printf("Coolio service ready at http://%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func (*context) contentHandler(w http.ResponseWriter, req *http.Request, ps router.Params) {
	view := ps.ByName("view")
	version := ps.ByName("version")
	file := ps.ByName("file")

	if version == "json" || version == "yaml" || version == "asset" || version == "js" {
		file = version
		version = ""
	}

	if v, err := NewViewModel(view, version); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	} else {
		switch file {
		case "json":
			w.Header().Add("content-type", "application/json")
			w.Write([]byte(v.JSON()))
		case "yaml":
			w.Header().Add("content-type", "text/plain")
			w.Write([]byte(v.YAML()))
		case "asset":
			t := GetAssetTemplate()
			t.Execute(w, v)
		case "js":
			t := GetJavascriptTemplate()
			t.Execute(w, v)
		default:
			t := GetEditorTemplate()
			t.Execute(w, v)
		}

		log.Println("view:", v.View, "v:", v.Version, "file:", file)
	}
}

func (*context) editorHandler(w http.ResponseWriter, req *http.Request, ps router.Params) {
	v, err := NewViewModel(ps.ByName("view"), ps.ByName("version"))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}

	t := GetEditorTemplate()
	t.Execute(w, v)
}

func (*context) saveHandler(w http.ResponseWriter, req *http.Request, ps router.Params) {
	// form reader....
	if err := req.ParseForm(); err != nil {
		// write error message.
	}

	v := NewViewModelFromRaw(
		ps.ByName("view"),
		req.FormValue("yaml"),
		req.FormValue("html"),
		req.FormValue("js"))

	// save view
	if _, preview := req.URL.Query()["preview"]; !preview {
		v.Save()
	}

	t := GetJavascriptTemplate()
	x, err := t.Bytes(v)

	if err != nil {
		log.Println("err:", err)
		return
	}

	preview := GetPreviewTemplate()
	preview.Execute(w, template.JS(x))
}

func (*context) indexHandler(w http.ResponseWriter, req *http.Request, ps router.Params) {
	w.Write([]byte("index.html"))
}
