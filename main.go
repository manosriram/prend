package main

import (
	"github.com/manosriram/prend/pkg/cmd"
)

func main() {
	cmd.Execute()
	// if len(os.Args) == 1 {
	// log.Fatal("no cmd provided")
	// }

	// conf, err := core.Init()
	// if err != nil {
	// log.Fatal("error initializing prend")
	// }

	// client := http.DefaultClient
	// githubSvc := github.NewGithubService(*client)

	// creds, err := githubSvc.GetGithubCreds()
	// if err != nil {
	// log.Fatal("error initializing github creds")
	// }

	// cmd := os.Args[1]
	// switch cmd {
	// case "source":
	// githubSvc.GetSources(conf, creds)

	// case "init":
	// return
	// }

}
