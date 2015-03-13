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
const proxy = "https://17afee27.ngrok.com"

func main() {
	router := httprouter.New()
	router.GET("/v1/:view", viewHandler)
	router.GET("/v1/:view/coolio-sp-asset.js", assetHandler)
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

func assetHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	view := ps.ByName("view")
	u := fmt.Sprintf("%s/v1/%s", proxy, view)
	fmt.Println("asset:", u)
	// fmt.Printf("%s://%s@%s/%s?%s#fragment", u.Scheme, u.User, u.Host, u.Path, u.Query().Encode())
	Sharepoint(w, u)
}
