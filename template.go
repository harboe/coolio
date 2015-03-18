package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	templateDir = "templates"
)

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

func GetAssetTemplate() Template {
	return Template{
		template.Must(template.ParseFiles(templateDir + "/asset.tmpl")),
		"application/javascript",
	}
}

func GetJavascriptTemplate() Template {
	return Template{template.Must(
		template.ParseFiles(
			templateDir+"/coolio.tmpl",
			templateDir+"/sharepoint.js",
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
	err := filepath.Walk(templateDir+"/"+name, func(path string, info os.FileInfo, err error) error {
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
