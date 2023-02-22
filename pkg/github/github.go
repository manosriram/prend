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

func getRawUrl(source data.Source, path string, token string) *data.RawUrl {
	comSplit := strings.Split(source.RepoUrl, ".git")
	gitSplit := strings.Split(comSplit[0], "/")
	username := gitSplit[3]
	repo := gitSplit[4]

	u := "https://api.github.com/repos/%s/%s/contents/%s?"
	var apiUrl string
	if source.Branch != "" {
		u += "ref=" + source.Branch + "&"
	}
	if token != "" {
		u += "token=" + token
	}
	apiUrl = fmt.Sprintf(u, username, repo, path)
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

func getRepoDefaultBranch(url string, creds *data.GithubCreds) (*data.GithubFileTree, error) {
	r, err := api.Request(url, creds.Token, http.MethodGet)
	if err != nil {
		return &data.GithubFileTree{}, err
	}
	var ff data.GithubFileTree
	err = json.Unmarshal([]byte(r), &ff)
	if err != nil {
		return &data.GithubFileTree{}, err
	}
	return &ff, nil
}

func GetRepoTree(source data.Source, creds *data.GithubCreds) (*data.RepoTree, error) {
	g := getRawUrl(source, "", creds.Token)

	u := "https://api.github.com/repos/manosriram-youtube/reddit_backend"
	fileTree, err := getRepoDefaultBranch(u, creds)
	if err != nil {
		fmt.Errorf("error getting default branch %s\n", err.Error())
	}

	var shaHashPath string
	if creds.Token != "" {
		shaHashPath = fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/%s?token=%s", g.Username, g.Repository, creds.Token, fileTree.DefaultBranch)
	} else {
		shaHashPath = fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/%s", g.Username, g.Repository, fileTree.DefaultBranch)
	}

	b, err := api.Request(shaHashPath, creds.Token, http.MethodPatch)
	if err != nil {
		fmt.Errorf("error getting api %s, %s\n", shaHashPath, err.Error())
	}

	var ff data.RepoTree
	err = json.Unmarshal([]byte(b), &ff)
	return &ff, nil
}

func GetFilesFromGithub(source data.Source, creds *data.GithubCreds) error {
	for _, path := range source.Paths {
		g := getRawUrl(source, path.SourcePath, creds.Token)
		protoFileList, err := api.Request(g.Url, creds.Token, http.MethodGet)
		if err != nil {
			// fmt.Printf("source %s not found. set GITHUB_TOKEN env var for private repos\n", (source.RepoUrl + "/" + path.SourcePath))
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

					fileContent, err := api.Request(entry.DownloadUrl, creds.Token, http.MethodGet)
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
