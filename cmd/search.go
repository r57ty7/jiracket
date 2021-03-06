/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/r57ty7/jiracket/infrastructure/jira"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runSearch,
}

var (
	project   string
	user      string
	parameter []string
)

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	searchCmd.Flags().StringVarP(&project, "project", "p", "", "Project for search")
	searchCmd.Flags().StringVarP(&user, "user", "u", "", "Assignee for search")
	searchCmd.Flags().StringArrayVar(&parameter, "parameter", []string{}, "Optional parameter [key=value]")
}

func runSearch(cmd *cobra.Command, args []string) {

	client, err := jira.NewClient(nil, config.Jira.BaseURL, config.Jira.Username, config.Jira.Password)
	if err != nil {
		cmd.PrintErrf("%v\n", err)
		os.Exit(1)
	}

	repo := jira.NewSearchRepository(client)

	param := make(map[string]string)
	param["project"] = project
	param["assignee"] = user

	jql := make([]string, 0, len(param))
	for k, v := range param {
		jql = append(jql, fmt.Sprintf("%v = %v", k, v))
	}

	jql = append(jql, parameter...)

	issues, err := repo.Search(context.Background(), strings.Join(jql, " AND "))
	if err != nil {
		cmd.PrintErrf("%v\n", err)
	}

	for _, v := range issues {
		fmt.Printf("%v\t%v\n", v.Key, v.Fields.Summary)
	}

}
