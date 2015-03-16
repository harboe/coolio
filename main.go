package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	router "github.com/julienschmidt/httprouter"
)

const port = "localhost:8080"
const proxy = "https://17afee27.ngrok.com"

type context struct{}

func main() {
	ctx := &context{}
	r := router.New()
	// content
	r.GET("/v1/:view/:version", ctx.contentHandler)
	r.GET("/v1/:view/:version/:file", ctx.contentHandler)

	// editor
	r.GET("/editor/:view", ctx.editorHandler)
	r.GET("/editor/:view/:version", ctx.editorHandler)
	// index
	r.GET("/", ctx.indexHandler)

	fmt.Printf("rest service ready at http://%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func (*context) contentHandler(w http.ResponseWriter, req *http.Request, ps router.Params) {
	file := strings.ToLower(ps.ByName("file"))

	if v, err := NewViewModel(ps.ByName("view"), ps.ByName("version")); err != nil {
		log.Println("err:", err)
		w.WriteHeader(404)
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
			w.Header().Add("content-type", "application/javascript")
			t := GetAssetTemplate()
			t.Execute(w, v)
		case "js":
			w.Header().Add("content-type", "application/javascript")
			t := GetJavascriptTemplate()
			t.Execute(w, v)
		default:
			w.Header().Add("content-type", "text/html")
			t := GetHtmlTemplate()
			t.Execute(w, v)
		}

		log.Println("view:", v.View, "v:", v.Version, "file:", file)
	}
}

func (*context) editorHandler(w http.ResponseWriter, req *http.Request, ps router.Params) {
	w.Write([]byte("editor.html"))
}

func (*context) indexHandler(w http.ResponseWriter, req *http.Request, ps router.Params) {
	w.Write([]byte("index.html"))
}
