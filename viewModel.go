package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	ViewModel struct {
		View       string
		Version    string
		yaml       template.HTML
		customHTML template.HTML
		customJS   template.HTML
	}
	HtmlTemplate struct {
		Name         string        `json:"name"`
		TemplateOnly bool          `json:"templateOnly"`
		Html         template.HTML `json:"html"`
	}
)

func NewViewModelFromRaw(view, yaml, html, js string) *ViewModel {
	return &ViewModel{
		View:       view,
		Version:    nextVersion(view),
		yaml:       template.HTML(yaml),
		customHTML: template.HTML(html),
		customJS:   template.HTML(js),
	}
}

func NewViewModel(view, version string) (*ViewModel, error) {
	if len(version) == 0 {
		version = latestVersion(view)
	}

	vm := &ViewModel{View: view, Version: version}

	if _, err := os.Stat(vm.ViewDir()); os.IsNotExist(err) {
		return nil, errors.New("not found")
	}

	return vm, nil
}

func latestVersion(view string) string {
	info, _ := ioutil.ReadDir("views/" + view)
	return info[len(info)-1].Name()
}

func nextVersion(view string) string {
	ver, _ := strconv.ParseInt(latestVersion(view), 10, 0)
	return fmt.Sprintf("%v", (ver + 1))
}

func (v *ViewModel) JSON() template.HTML {
	b, err := json.Marshal(v.Group())

	if err != nil {
		log.Fatal(err)
	}

	return template.HTML(b)
}

func (v *ViewModel) YAML() template.HTML {
	if len(v.yaml) > 0 {
		return v.yaml
	}

	file := v.ViewDir() + "/layout.yaml"
	return template.HTML(readFile(file))
}

func (v *ViewModel) CustomHTML() HTML {
	var b []byte

	if len(v.customHTML) > 0 {
		b = []byte(v.customHTML)
	} else {
		file := v.ViewDir() + "/custom.html"
		b = []byte(readFile(file))
	}

	return HTML(b)
}

func (v *ViewModel) CustomJS() template.HTML {
	if len(v.customJS) > 0 {
		return v.customJS
	}
	file := v.ViewDir() + "/custom.js"
	return template.HTML(readFile(file))
}

func readFile(file string) []byte {
	if b, err := ioutil.ReadFile(file); err == nil {
		return b
	}
	return []byte{}
}

func (v *ViewModel) Group() (g Group) {
	if len(v.yaml) > 0 {
		err := yaml.Unmarshal([]byte(v.yaml), &g)

		if err != nil {
			log.Println(err)
		}
	} else if b := v.YAML(); len(b) > 0 {
		err := yaml.Unmarshal([]byte(b), &g)

		if err != nil {
			log.Println(err)
		}
	}

	return g
}

func (v *ViewModel) Overrides() template.HTML {
	b, err := json.Marshal(v.Group().AllParameters())

	if err != nil {
		log.Fatal(err)
	}

	return template.HTML(b)
}

func (v *ViewModel) Templates() template.HTML {
	list := []HtmlTemplate{}
	traverseTemplates("html", ".html", func(path string, b []byte) {
		filename := filepath.Base(path)
		name := filename[:strings.LastIndex(filename, ".")]

		tmpl := HtmlTemplate{
			Name:         name,
			TemplateOnly: strings.HasPrefix(string(name), "coolio-"),
			Html:         template.HTML(HTML(b).Inline()),
		}

		list = append(list, tmpl)
	})

	b, err := json.Marshal(list)

	if err != nil {
		log.Fatal(err)
	}

	return template.HTML(b)
}

func (v *ViewModel) Javascript() template.HTML {
	buf := bytes.Buffer{}
	traverseTemplates("js", ".js", func(path string, b []byte) {
		buf.Write(b)
	})
	return template.HTML(buf.Bytes())
}

func (v *ViewModel) JsLibraries() template.HTML {
	buf := bytes.Buffer{}
	traverseTemplates("libs", ".js", func(path string, b []byte) {
		buf.Write(b)
	})
	return template.HTML(buf.Bytes())
}

func (v *ViewModel) String() string {
	return fmt.Sprintf(proxy+"/v1/%s/%s", v.View, v.Version)
}

func (v *ViewModel) ViewDir() string {
	return fmt.Sprintf("views/%s/%s", v.View, v.Version)
}

func (v *ViewModel) Sum() string {
	h := md5.New()
	h.Write([]byte(v.YAML()))
	h.Write([]byte(v.CustomHTML()))
	h.Write([]byte(v.CustomJS()))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (v *ViewModel) HasChanges() bool {
	file := fmt.Sprintf("views/%s/%s/md5", v.View, latestVersion(v.View))
	if b := readFile(file); len(b) > 0 && v.Sum() == string(b) {
		return false
	}
	return true
}

func (v *ViewModel) Save() error {
	dir := v.ViewDir()

	if ok := v.HasChanges(); !ok {
		log.Println("ok?", ok)
		return errors.New("allready saved")
	}

	// first create the directory
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		return err
	}
	// save yaml, json, custom html & js
	ioutil.WriteFile(dir+"/layout.yaml", []byte(v.yaml), os.ModePerm)
	ioutil.WriteFile(dir+"/layout.json", []byte(v.JSON()), os.ModePerm)
	ioutil.WriteFile(dir+"/custom.html", []byte(v.customHTML), os.ModePerm)
	ioutil.WriteFile(dir+"/custom.js", []byte(v.customJS), os.ModePerm)

	ioutil.WriteFile(dir+"/md5", []byte(v.Sum()), os.ModePerm)

	// save sharepoint asset file
	asset := GetAssetTemplate()
	if b, err := asset.Bytes(v); err != nil {
		return err
	} else {
		ioutil.WriteFile(dir+"/asset.js", b, os.ModePerm)
	}

	// save javascript file
	js := GetJavascriptTemplate()
	if b, err := js.Bytes(v); err != nil {
		return err
	} else {
		ioutil.WriteFile(dir+"/coolio.js", b, os.ModePerm)
	}

	return nil
}
