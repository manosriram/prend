package data

import "time"

type Path struct {
	SourcePath      string `yaml:"source_proto_path"`
	DestinationPath string `yaml:"destination_proto_path"`
}

type Source struct {
	// Destination_path string   `yaml:"destination_path"`
	RepoUrl string `yaml:"repo_url"`
	Branch  string `yaml: "branch"`
	Paths   []Path `yaml:"paths"`
}

type Conf struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Sources     []Source `yaml: "sources"`
}

type LockFile struct {
	Repo        string    `yaml:"repo"`
	Commit      string    `yaml:"commit"`
	Branch      string    `yaml:"branch"`
	LastUpdated time.Time `yaml:last_updated"`
}

type LockYamlData struct {
	Sources []LockFile `yaml:"sources"`
}
