package core

import (
	"io/ioutil"
	"log"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type service struct {
	logger *zap.SugaredLogger
}

func NewVendorService(logger *zap.SugaredLogger) *service {
	return &service{
		logger: logger,
	}
}

type Source struct {
	// Destination_path string   `yaml:"destination_path"`
	RepoUrl         string `yaml:"repo_url"`
	Branch          string `yaml: "branch"`
	DestinationPath string `yaml:"destination_proto_path"`
	SourcePath      string `yaml:"source_proto_path"`
}

type Conf struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Sources     []Source `yaml: "sources"`
}

func loadYaml() (*Conf, error) {
	yamlFile, err := ioutil.ReadFile("prend.yaml")

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
