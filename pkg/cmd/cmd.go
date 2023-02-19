package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"prend/pkg/core"
	github "prend/pkg/github_auth"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "prend",
	Short: "protocol buffer vendor",
	Long:  "vendoring for protobufs",
}

var sourceCmd = &cobra.Command{
	Use:   "source",
	Short: "update all the sources from prend.config file",
	Long:  "pulls all the updated proto files from the given sources",
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
	rootCmd.AddCommand(sourceCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
