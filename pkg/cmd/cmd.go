package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	github "github.com/manosriram/prend/pkg/github"

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

		// remove "vendor" folder
		err = os.RemoveAll("vendors")
		if err != nil {
			log.Fatal("error removing vendors directory during initialization")
		}

		core.GetSources(conf, creds)
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "cleans up all source_proto_path directories",
	Long:  "cleans up source_proto_path directories from prend.yaml file",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := core.Init()
		if err != nil {
			log.Fatal("error initializing prend")
		}
		core.Clean(conf)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(cleanCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
