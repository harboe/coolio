package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/yaml.v2"
)

const port = "localhost:8080"

func main() {
	router := httprouter.New()
	router.GET("/v1/:view", viewHandler)
	router.GET("/v1/:view/debug", debugHandler)

	fmt.Printf("rest service ready at http://%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func viewHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	view := ps.ByName("view")
	fmt.Println("view:", view)
	b, err := ioutil.ReadFile("views/" + view + ".yaml")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var g Group

	if err := yaml.Unmarshal(b, &g); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/javascript")
	GenerateJS(w, g)
}

func debugHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	view := ps.ByName("view")
	fmt.Println("debug:", view)
	Debug(w, "/v1/"+view)
}
