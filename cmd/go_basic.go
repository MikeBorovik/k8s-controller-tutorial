/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// goBasicCmd represents the go-basic command
var goBasicCmd = &cobra.Command{
	Use:   "go-basic",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Go basic code to run functions
		k8s := Kubernetes{
			Name:    "k8s-demo-cluster",
			Version: "1.31",
			Users:   []string{"alex", "den"},
			NodeNumber: func() int {
				return 10
			},
		}

		//print users
		k8s.GetUsers()

		//add new user to struct
		k8s.AddNewUser("anonymous")

		//print users one more time
		k8s.GetUsers()
	},
}

func init() {
	rootCmd.AddCommand(goBasicCmd)
}

// My go basic fucntions here
type Kubernetes struct {
	Name       string     `json:"name"`
	Version    string     `json:"version"`
	Users      []string   `json:"users,omitempty"`
	NodeNumber func() int `json:"-"`
}

func (k8x Kubernetes) GetUsers() {
	for _, user := range k8x.Users {
		fmt.Println(user)
	}
}

func (k9s *Kubernetes) AddNewUser(user string) {
	k9s.Users = append(k9s.Users, user)
}
