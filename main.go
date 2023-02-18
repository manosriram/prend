package main

import (
	"prend/pkg/core"
	github "prend/pkg/github_auth"
)

func main() {
	conf, _ := core.Init()
	// fmt.Println(conf)

	creds, _ := github.GetGithubCreds()
	github.GetSources(conf, creds)

	// cmd := os.Args[1]

	// fmt.Println(cmd)

	// switch cmd {

	// }

}
