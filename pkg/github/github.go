package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/manosriram/prend/pkg/api"
	"github.com/manosriram/prend/pkg/data"
)

type service struct {
	client http.Client
}

func NewGithubService(client http.Client) *service {
	return &service{
		client: client,
	}
}

func getRawUrl(url string, path string, token string) *data.RawUrl {
	comSplit := strings.Split(url, ".git")
	gitSplit := strings.Split(comSplit[0], "/")
	username := gitSplit[3]
	repo := gitSplit[4]

	u := "https://api.github.com/repos/%s/%s/contents/%s"
	var apiUrl string
	if token != "" {
		u += "?token=%s"
		apiUrl = fmt.Sprintf(u, username, repo, path, token)
	} else {
		apiUrl = fmt.Sprintf(u, username, repo, path)
	}
	return &data.RawUrl{
		Username:   username,
		Repository: repo,
		Url:        apiUrl,
	}
}

func (svc *service) GetGithubCreds() (*data.GithubCreds, error) {
	token := os.Getenv("GITHUB_TOKEN")

	return &data.GithubCreds{
		Token: token,
	}, nil
}

func GetRepoTree(source data.Source, creds *data.GithubCreds) (*data.RepoTree, error) {
	g := getRawUrl(source.RepoUrl, "", creds.Token)
	var shaHashPath string
	if creds.Token != "" {
		shaHashPath = fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/master?token=%s", g.Username, g.Repository, creds.Token)
	} else {
		shaHashPath = fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/master", g.Username, g.Repository)
	}
	b, err := api.Get(shaHashPath, creds.Token)
	if err != nil {
		fmt.Errorf("error getting api %s, %s", shaHashPath, err.Error())
	}

	var ff data.RepoTree
	err = json.Unmarshal([]byte(b), &ff)
	return &ff, nil
}

func GetFilesFromGithub(source data.Source, creds *data.GithubCreds) error {
	for _, path := range source.Paths {
		g := getRawUrl(source.RepoUrl, path.SourcePath, creds.Token)
		protoFileList, err := api.Get(g.Url, creds.Token)
		if err != nil {
			fmt.Printf("source %s not found. set GITHUB_TOKEN env var for private repos\n", (source.RepoUrl + "/" + path.SourcePath))
			return err
		}

		if _, err := os.Stat(path.DestinationPath); os.IsNotExist(err) {
			os.MkdirAll(path.DestinationPath, os.ModeDir|0755)
		}

		var r []data.GithubFileTree

		err = json.Unmarshal([]byte(protoFileList), &r)
		if err != nil {
			return err
		}

		for _, entry := range r {
			if entry.Type == "file" {
				ext := strings.Split(entry.Name, ".")
				if len(ext) > 1 && ext[1] == "proto" {
					fmt.Printf("source %s found\n", entry.HtmlUrl)

					fileContent, err := api.Get(entry.DownloadUrl, creds.Token)
					if err != nil {
						return err
					}

					vendorPath := fmt.Sprintf("vendor")

					if path.DestinationPath != "" {
						vendorPath = path.DestinationPath
					}

					f, err := os.OpenFile(fmt.Sprintf("%s/%s", vendorPath, entry.Name), os.O_RDWR|os.O_CREATE, 0755)
					if err != nil {
						return err
					}
					f.Write([]byte(fileContent))
				}
			}
		}
	}

	return nil
}
