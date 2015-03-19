package main

import (
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
	TemplateKey string
)

func (key TemplateKey) HasViewModel() bool {
	return !strings.HasPrefix(string(key), "coolio-")
}

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

func (v *ViewModel) Parameters() []Parameter {
	return v.Group().AllParameters()
}

func (v *ViewModel) Bundles() []template.HTML {
	list := []template.HTML{}
	traverseTemplates("bundles", ".js", func(path string, b []byte) {
		list = append(list, template.HTML(b))
	})
	return list
}

func (v *ViewModel) Templates() map[TemplateKey]HTML {
	dic := map[TemplateKey]HTML{}
	traverseTemplates("html", ".html", func(path string, b []byte) {
		filename := filepath.Base(path)
		templateName := filename[:strings.LastIndex(filename, ".")]

		// tmpl := removeWhitespace(b)
		// tmpl = strings.Replace(tmpl, "'", "\\'", -1)

		dic[TemplateKey(templateName)] = HTML(b)
	})
	return dic
}

func (v *ViewModel) Javascript() []template.HTML {
	list := []template.HTML{}
	traverseTemplates("js", ".js", func(path string, b []byte) {
		list = append(list, template.HTML(b))
	})
	return list
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

// func removeWhitespace(b []byte) string {
// 	tmpl := strings.Replace(string(b), "\n", "", -1)
// 	tmpl = strings.Replace(tmpl, "\r", "", -1)
// 	return strings.Replace(tmpl, "\t", "", -1)
// }
