package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type HTML template.HTML

func (h HTML) Inline() template.HTML {
	tmpl := strings.Replace(string(h), "\n", "", -1)
	tmpl = strings.Replace(tmpl, "\r", "", -1)
	tmpl = strings.Replace(tmpl, "\t", "", -1)
	// tmpl = strings.Replace(tmpl, "'", "\\'", -1)

	return template.HTML(tmpl)
}

type Template struct {
	*template.Template
	contentType string
}

func (t Template) Execute(w http.ResponseWriter, v interface{}) {
	w.Header().Add("content-type", t.contentType)
	w.WriteHeader(http.StatusOK)

	if err := t.Template.Execute(w, v); err != nil {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (t Template) Bytes(v interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	err := t.Template.Execute(buf, v)
	return buf.Bytes(), err
}

// // First we create a FuncMap with which to register the function.
var funcMap = template.FuncMap{
	"last": func(a, b interface{}) (interface{}, error) {
		return a.(int) == b.(int), nil
	},
}

func GetAssetTemplate() Template {
	return Template{
		template.Must(template.ParseFiles(templateDir + "/asset.tmpl")),
		"application/javascript",
	}
}

func GetJavascriptTemplate() Template {
	return Template{template.Must(
		template.ParseFiles(
			templateDir + "/coolio.js",
		)),
		"application/javascript",
	}
}

func GetEditorTemplate() Template {
	return Template{
		template.Must(template.ParseFiles(templateDir + "/editor.tmpl")),
		"text/html",
	}
}

func GetPreviewTemplate() Template {
	return Template{
		template.Must(template.ParseFiles(templateDir + "/preview.tmpl")),
		"text/html",
	}
}

func traverseTemplates(name, ext string, op func(path string, b []byte)) {
	log.Println(templateDir, name, ext)
	err := filepath.Walk(templateDir+"/"+name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || (len(ext) > 0 && filepath.Ext(path) != ext) {
			return nil
		}

		b, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		op(path, b)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
