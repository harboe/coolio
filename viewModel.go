package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type ViewModel struct {
	View    string
	Version string
	body    []byte
}

func NewViewModelFromRaw(b []byte) (*ViewModel, error) {
	return &ViewModel{body: b}, nil
}

func NewViewModel(view, version string) (*ViewModel, error) {
	if len(version) == 0 {
		info, _ := ioutil.ReadDir("views/" + view)
		version = info[len(info)-1].Name()
	}

	vm := &ViewModel{view, version, []byte{}}

	if _, err := os.Stat(vm.ViewDir()); os.IsNotExist(err) {
		return nil, errors.New("not found")
	}

	return vm, nil
}

func (v *ViewModel) JSON() template.HTML {
	b, err := json.Marshal(v.Group())

	if err != nil {
		log.Fatal(err)
	}

	return template.HTML(b)
}

func (v *ViewModel) YAML() template.HTML {
	file := v.ViewDir() + "/layout.yaml"

	if b, err := ioutil.ReadFile(file); err != nil {
		log.Println("err")
	} else {
		return template.HTML(b)
	}

	return ""
}

func (v *ViewModel) Group() (g Group) {
	if len(v.body) > 0 {
		err := yaml.Unmarshal(v.body, &g)

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
	traverseTemplates("bundles", "", func(path string, b []byte) {
		list = append(list, template.HTML(b))
	})
	return list
}

func (v *ViewModel) Templates() map[string]template.HTML {
	dic := map[string]template.HTML{}
	traverseTemplates("html", ".html", func(path string, b []byte) {
		filename := filepath.Base(path)
		templateName := filename[:strings.LastIndex(filename, ".")]

		tmpl := strings.Replace(string(b), "\n", "", -1)
		tmpl = strings.Replace(tmpl, "'", "\\'", -1)

		dic[templateName] = template.HTML(tmpl)
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
