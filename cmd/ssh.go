// Copyright © 2017 Manuel Gauto <github.com/twa16>
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
	"fmt"

	"github.com/spf13/cobra"
	"os"
	"bufio"
	"strings"
	"strconv"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		_, err := GetSavedSession()
		if err != nil {
			fmt.Println("Error Getting Sessions: "+err.Error())
			os.Exit(1)
		}
		//Print spaces
		spaces := GetSpaces()
		fmt.Println("Your Spaces:")
		for i, space := range spaces {
			fmt.Printf("[%d] %s\n", i, space.FriendlyName)
		}

		//Get Space Index
		fmt.Print("What Space would you like to connect to? ")
		spaceIndexString, err := reader.ReadString('\n')
		spaceIndexString = strings.TrimSpace(spaceIndexString)
		spaceIndex, err := strconv.Atoi(spaceIndexString)
		if err != nil {
			fmt.Println("Invalid Choice: "+err.Error())
			os.Exit(1)
		}
		if spaceIndex < 0 || spaceIndex >= len(spaces) {
			fmt.Println("Invalid Choice: Space does not exist")
			os.Exit(1)
		}
		
	},
}

func init() {
	RootCmd.AddCommand(sshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
