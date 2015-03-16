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
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type ViewModel struct {
	View    string
	Version uint
}

func NewViewModel(view, version string) (*ViewModel, error) {
	v, _ := strconv.ParseUint(version, 10, 0)
	vm := &ViewModel{view, uint(v)}

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
	if b := v.YAML(); len(b) > 0 {
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

func (v *ViewModel) Templates() map[string]template.HTML {
	dic := map[string]template.HTML{}
	traverseTemplates("html", func(path string, b []byte) {
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
	traverseTemplates("js", func(path string, b []byte) {
		list = append(list, template.HTML(b))
	})
	return list
}

func (v *ViewModel) String() string {
	return fmt.Sprintf("/v1/%s/%v/js", v.View, v.Version)
}

func (v *ViewModel) ViewDir() string {
	return fmt.Sprintf("%s/%s/%v", viewDir, v.View, v.Version)
}
