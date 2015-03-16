package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	router "github.com/julienschmidt/httprouter"
)

const port = "localhost:8080"
const proxy = "http://localhost:8080"

type context struct{}

func main() {
	ctx := &context{}

	r := router.New()
	// content
	r.GET("/v1/:view/:version", ctx.contentHandler)
	r.GET("/v1/:view/:version/:file", ctx.contentHandler)

	// preview
	r.POST("/preview", ctx.previewHandler)

	// editor
	r.GET("/editor/:view", ctx.editorHandler)
	r.GET("/editor/:view/:version", ctx.editorHandler)
	// index
	r.GET("/", ctx.indexHandler)

	r.ServeFiles("/static/*filepath", http.Dir("static"))

	fmt.Printf("rest service ready at http://%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func (*context) contentHandler(w http.ResponseWriter, req *http.Request, ps router.Params) {
	view := ps.ByName("view")
	version := ps.ByName("version")
	file := ps.ByName("file")

	if len(ps) == 2 {
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
		case "editor":
			t := GetEditorTemplate()
			t.Execute(w, v)
		default:
			t := GetHtmlTemplate()
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

func (*context) previewHandler(w http.ResponseWriter, req *http.Request, ps router.Params) {
	defer req.Body.Close()
	b, _ := ioutil.ReadAll(req.Body)
	v, _ := NewViewModelFromRaw(b)
	t := GetJavascriptTemplate()
	x, _ := t.ExecuteBytes(v)

	preview := GetPreviewTemplate()
	preview.Execute(w, template.JS(x))
}

func (*context) indexHandler(w http.ResponseWriter, req *http.Request, ps router.Params) {
	w.Write([]byte("index.html"))
}
