package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	templateDir = "templates"
	viewDir     = "views"
)

func GetAssetTemplate() *template.Template {
	return template.Must(template.ParseFiles(templateDir + "/asset.tmpl"))
}

func GetJavascriptTemplate() *template.Template {
	return template.Must(template.ParseFiles(
		templateDir+"/coolio.tmpl",
		templateDir+"/sharepoint.js",
	))
}

func GetHtmlTemplate() *template.Template {
	return template.Must(template.ParseFiles(templateDir + "/debug.tmpl"))
}

func traverseTemplates(name string, op func(path string, b []byte)) {
	err := filepath.Walk(templateDir+"/"+name, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != "."+name {
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
