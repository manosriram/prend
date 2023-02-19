package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"prend/pkg/core"
	"strings"
)

// {
// "name": "Base.hpp",
// "path": "Bloom-Filter/Base.hpp",
// "sha": "2fdc98c7db1b4de13adb0a89f1403964e22a1479",
// "size": 281,
// "url": "https://api.github.com/repos/manosriram/Data-Structures/contents/Bloom-Filter/Base.hpp?ref=master",
// "html_url": "https://github.com/manosriram/Data-Structures/blob/master/Bloom-Filter/Base.hpp",
// "git_url": "https://api.github.com/repos/manosriram/Data-Structures/git/blobs/2fdc98c7db1b4de13adb0a89f1403964e22a1479",
// "download_url": "https://raw.githubusercontent.com/manosriram/Data-Structures/master/Bloom-Filter/Base.hpp",
// "type": "file",
// "_links": {
// "self": "https://api.github.com/repos/manosriram/Data-Structures/contents/Bloom-Filter/Base.hpp?ref=master",
// "git": "https://api.github.com/repos/manosriram/Data-Structures/git/blobs/2fdc98c7db1b4de13adb0a89f1403964e22a1479",
// "html": "https://github.com/manosriram/Data-Structures/blob/master/Bloom-Filter/Base.hpp"
// }
// },

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
	Client http.Client
}

type GithubCreds struct {
	Token string
}

func getRawUrl(url string, path string, token string) string {
	comSplit := strings.Split(url, ".git")
	gitSplit := strings.Split(comSplit[0], "/")
	username := gitSplit[3]
	repo := gitSplit[4]

	u := "https://api.github.com/repos/%s/%s/contents/%s"
	var x string
	if token != "" {
		u += "?token=%s"
		x = fmt.Sprintf(u, username, repo, path, token)
	} else {
		x = fmt.Sprintf(u, username, repo, path)
	}
	return x
}

func GetGithubCreds() (*GithubCreds, error) {
	token := os.Getenv("GITHUB_TOKEN")

	return &GithubCreds{
		Token: token,
	}, nil
}

func get(url string) (string, error) {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("source %s not found", url))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getFilesFromGithub(source core.Source, creds *GithubCreds) {
	protoUrl := getRawUrl(source.RepoUrl, source.SourcePath, creds.Token)
	protoFileList, err := get(protoUrl)
	if err != nil {
		fmt.Printf("source %s not found\n", source.RepoUrl)
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

				fileContent, _ := get(entry.DownloadUrl)
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

func GetSources(conf *core.Conf, creds *GithubCreds) {
	for _, source := range conf.Sources {
		getFilesFromGithub(source, creds)
		// if source.Branch == "" {
		// fmt.Printf("missing branch for %s, defaulting to master\n", source.Repo_url)
		// source.Branch = "master"
		// }

		// fmt.Println(source)
	}
}
