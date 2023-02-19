package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	github "github.com/manosriram/prend/pkg/github_auth"

	"github.com/manosriram/prend/pkg/core"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "prend",
	Short: "protocol buffer vendor",
	Long:  "vendoring for protobufs",
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch all the sources from prend.config file",
	Long:  "fetches all the updated proto files from the given sources",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := core.Init()
		if err != nil {
			log.Fatal("error initializing prend")
		}

		client := http.DefaultClient
		githubSvc := github.NewGithubService(*client)
		creds, err := githubSvc.GetGithubCreds()
		if err != nil {
			log.Fatal("error initializing github creds")
		}

		githubSvc.GetSources(conf, creds)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
