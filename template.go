package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const templateDir = "templates"

// var coolio *template.Template
type coolioViewModel struct {
	Templates  map[string]template.HTML
	Javascript []template.HTML
	Layout     template.HTML
	Parameters []Parameter
}

func loadCoolio() (*template.Template, coolioViewModel) {
	var err error
	var templates map[string]template.HTML
	var javascript []template.HTML

	if templates, err = loadHtml(templateDir + "/html"); err != nil {
		log.Fatal(err)
	}
	if javascript, err = traverseDir(templateDir + "/js"); err != nil {
		log.Fatal(err)
	}

	vm := coolioViewModel{
		Templates:  templates,
		Javascript: javascript,
	}

	return template.Must(template.ParseFiles(templateDir + "/coolio.tmpl")), vm
}

func GenerateJS(w io.Writer, g Group) error {
	js, err := json.Marshal(g)

	if err != nil {
		return err
	}

	coolio, vm := loadCoolio()
	vm.Layout = template.HTML(js)
	vm.Parameters = g.AllParameters()

	return coolio.Execute(w, vm)
}

func Sharepoint(w io.Writer, view string) error {
	t := template.Must(template.ParseFiles(templateDir + "/sharepoint.tmpl"))
	return t.Execute(w, view)
}

func Debug(w io.Writer, view string) error {
	t := template.Must(template.ParseFiles(templateDir + "/debug.tmpl"))
	return t.Execute(w, view)
}

func traverseDir(dir string) (list []template.HTML, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		js, err := loadSingleFile(path)

		if err != nil {
			return err
		}

		list = append(list, js)
		return nil
	})

	return list, err
}

func loadSingleFile(file string) (template.HTML, error) {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		return "", err
	}

	return template.HTML(b), nil
}

func loadHtml(dir string) (dic map[string]template.HTML, err error) {
	dic = map[string]template.HTML{}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}

		filename := strings.Replace(path, dir+"/", "", 1)
		templateName := filename[:strings.LastIndex(filename, ".")]

		tmpl := strings.Replace(string(b), "\n", "", -1)
		tmpl = strings.Replace(tmpl, "'", "\\'", -1)

		dic[templateName] = template.HTML(tmpl)
		return nil
	})

	return dic, err
}
