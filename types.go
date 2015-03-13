package main

type (
	Group struct {
		Parameter  `yaml:",inline"`
		Groups     []Group     `yaml:"groups,flow" json:"groups"`
		Parameters []Parameter `yaml:"params,flow" json:"params"`
	}
	Parameter struct {
		Id          string            `yaml:"id" json:"id"`
		Name        string            `yaml:"name,omitempty" json:"name,omitempty"`
		Description string            `yaml:"desc,omitempty" json:"desc,omitempty"`
		Type        string            `yaml:"type" json:"type"`
		Help        string            `yaml:"help,omitempty" json:"help,omitempty"`
		Aux         map[string]string `yaml:",flow,omitempty" json:"aux,omitempty"`
	}
)
