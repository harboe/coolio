package main

type (
	Group struct {
		Parameter  `yaml:",inline"`
		Groups     []Group     `yaml:"groups,omitempty" json:"groups,omitempty"`
		Parameters []Parameter `yaml:"params,omitempty" json:"params,omitempty"`
	}
	Parameter struct {
		Id          string            `yaml:"id,omitempty" json:"id,omitempty"`
		Name        string            `yaml:"name,omitempty" json:"name,omitempty"`
		Description string            `yaml:"desc,omitempty" json:"desc,omitempty"`
		Type        string            `yaml:"type,omitempty" json:"type,omitempty"`
		Help        string            `yaml:"help,omitempty" json:"help,omitempty"`
		Aux         map[string]string `yaml:",flow,omitempty" json:"aux,omitempty"`
	}
)

func (g Group) AllParameters() (p []Parameter) {
	recursive(&p, g)
	return p
}

func recursive(list *[]Parameter, g Group) {
	for _, p := range g.Parameters {
		if len(p.Id) > 0 {
			*list = append(*list, p)
		}
	}

	for _, val := range g.Groups {
		recursive(list, val)
	}
}
