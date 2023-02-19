package main

import (
	"os"
	"prend/pkg/core"
	github "prend/pkg/github_auth"
)

func main() {
	conf, _ := core.Init()
	creds, _ := github.GetGithubCreds()

	cmd := os.Args[1]

	switch cmd {
	case "source":
		github.GetSources(conf, creds)

	case "init":
		return
	}

}
