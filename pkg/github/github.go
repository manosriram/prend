package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/manosriram/prend/pkg/api"
	"github.com/manosriram/prend/pkg/core"
)

type githubFolderTree struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        int32  `json:"size"`
	Url         string `json:"url"`
	DownloadUrl string `json:"download_url"`
	HtmlUrl     string `json:"html_url"`
	Type        string `json:"type"`
}

type service struct {
	client http.Client
}

func NewGithubService(client http.Client) *service {
	return &service{
		client: client,
	}
}

type GithubCreds struct {
	Token string
}

func (svc *service) getRawUrl(url string, path string, token string) string {
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
	return apiUrl
}

func (svc *service) GetGithubCreds() (*GithubCreds, error) {
	token := os.Getenv("GITHUB_TOKEN")

	return &GithubCreds{
		Token: token,
	}, nil
}

func (svc *service) getFilesFromGithub(source core.Source, creds *GithubCreds) {
	protoUrl := svc.getRawUrl(source.RepoUrl, source.SourcePath, creds.Token)
	protoFileList, err := api.Get(protoUrl)
	if err != nil {
		fmt.Printf("source %s not found. set GITHUB_TOKEN env var for private repos\n", (source.RepoUrl + "/" + source.SourcePath))
		return
	}

	if _, err := os.Stat(source.DestinationPath); os.IsNotExist(err) {
		os.MkdirAll(source.DestinationPath, os.ModeDir|0755)
	}

	var r []githubFolderTree
	err = json.Unmarshal([]byte(protoFileList), &r)
	if err != nil {
		panic(err)
	}

	for _, entry := range r {
		if entry.Type == "file" {
			ext := strings.Split(entry.Name, ".")
			if len(ext) > 1 && ext[1] == "proto" {
				fmt.Printf("source %s found\n", entry.HtmlUrl)

				fileContent, _ := api.Get(entry.DownloadUrl)
				path := fmt.Sprintf("vendor")

				if source.DestinationPath != "" {
					path = source.DestinationPath
				}

				f, err := os.OpenFile(fmt.Sprintf("%s/%s", path, entry.Name), os.O_RDWR|os.O_CREATE, 0755)
				if err != nil {
					fmt.Println(err)
				}
				f.Write([]byte(fileContent))
			}

		}
	}
}

func (svc *service) GetSources(conf *core.Conf, creds *GithubCreds) {
	for _, source := range conf.Sources {
		svc.getFilesFromGithub(source, creds)
	}
}
