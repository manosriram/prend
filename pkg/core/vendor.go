package core

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type source struct {
	Destination_path string   `yaml:"destination_path"`
	Repo_url         string   `yaml:"repo_url"`
	Proto_files      []string `yaml: "proto_files"`
	Branch           string   `yaml: "branch"`
}

type Conf struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Sources     []source `yaml: "sources"`
}

func loadYaml() (*Conf, error) {
	yamlFile, err := ioutil.ReadFile("dev.yaml")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	d := Conf{}
	err = yaml.Unmarshal(yamlFile, &d)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &d, nil
}

func Init() {
	loadYaml()
}
