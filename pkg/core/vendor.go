package core

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Source struct {
	// Destination_path string   `yaml:"destination_path"`
	Repo_url         string   `yaml:"repo_url"`
	Proto_file_paths []string `yaml: "proto_file_urls"`
	Branch           string   `yaml: "branch"`
	DestinationPath  string   `yaml:"dest_path"`
}

type Conf struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Sources     []Source `yaml: "sources"`
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

func Init() (*Conf, error) {
	return loadYaml()
}

func Update() {}

func updateSources() {}
