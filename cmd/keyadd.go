// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// keyaddCmd represents the keyadd command
var keyaddCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Default key path
		keyDefaultPath := "~/.ssh/id_rsa.pub"

		//Banner
		fmt.Println("Adding an SSH Public Key to your account")
		//Setup reader
		reader := bufio.NewReader(os.Stdin)
		//Get path to key
		fmt.Printf("Path to Public Key [%s]: ", keyDefaultPath)
		keyPathString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Failed to get input: " + err.Error())
			os.Exit(1)
		}
		//Cleanup string
		keyPathString = strings.TrimSpace(keyPathString)

		if keyPathString == "" {
			fmt.Printf("Assuming Default: %s\n", keyDefaultPath)
			keyPathString = keyDefaultPath
		}
	},
}

func init() {
	keyCmd.AddCommand(keyaddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keyaddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keyaddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
