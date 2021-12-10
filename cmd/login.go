/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"fmt"

	"github.com/bachelor-thesis-hown3d/chat-client/pkg/login"
	"github.com/spf13/cobra"
)

var (
	clientSecret string
	clientID     string
	issuerUrl    string
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to OAuth Provider",
	Long: `A longer description that spans multiple lines and likely contains examples

Example:

chat-client login --client-id=kubernetes --client-secret=kubernetes https://keycloak:8443/auth/realms/kubernetes
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := login.Login(cmd.Context(), clientID, clientSecret, args[0])
		if err != nil {
			fmt.Println(err)
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	f := loginCmd.Flags()
	f.StringVar(&clientSecret, "client-secret", "kubernetes", "Client Secret for OAuth")
	f.StringVar(&clientID, "client-id", "kubernetes", "Client ID for OAuth")
	//cobra.MarkFlagRequired(f, "client-secret")
	//cobra.MarkFlagRequired(f, "client-id")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
