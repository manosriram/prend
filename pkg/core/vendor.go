package core

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/manosriram/prend/pkg/data"
	"github.com/manosriram/prend/pkg/github"
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

func loadYaml() (*data.Conf, error) {
	yamlFile, err := ioutil.ReadFile("prend.yaml")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	d := data.Conf{}
	err = yaml.Unmarshal(yamlFile, &d)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &d, nil
}

func Init() (*data.Conf, error) {
	return loadYaml()
}

func writeLockFile(sources data.LockYamlData, f *os.File) error {
	yamlData, _ := yaml.Marshal(&sources)
	_, err := f.WriteString(string(yamlData))
	if err != nil {
		return err
	}

	return nil
}

func getSourcesFromGithub(conf *data.Conf, creds *data.GithubCreds) ([]*data.RepoTree, error) {
	fileName := "prend-lock.yaml"
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	f.Truncate(0)
	var sources data.LockYamlData

	var trees []*data.RepoTree
	for _, source := range conf.Sources {
		tree, err := github.GetRepoTree(source, creds)
		if err != nil {
			return nil, err
		}

		err = github.GetFilesFromGithub(source, creds)
		if err != nil {
			return trees, err
		}

		trees = append(trees, tree)
		sources.Sources = append(sources.Sources, data.LockFile{
			Repo:        source.RepoUrl,
			Branch:      "master",
			Commit:      tree.Commit.Sha,
			LastUpdated: time.Now(),
		})
	}
	err = writeLockFile(sources, f)
	if err != nil {
		return nil, err
	}

	f.Close()
	return trees, nil
}

func GetSources(conf *data.Conf, creds *data.GithubCreds) {
	_, err := getSourcesFromGithub(conf, creds)
	if err != nil {
		log.Fatal("error occurred ", err)
	}
}
