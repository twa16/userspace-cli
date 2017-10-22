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

	"bufio"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a space",
	Long:  `Deletes a space and all associated data.`,
	Run: func(cmd *cobra.Command, args []string) {
		execDeleteCommand(cmd, args)
	},
}

func execDeleteCommand(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)
	session, err := GetSavedSession()
	if err != nil {
		fmt.Println("Error Getting Sessions: " + err.Error())
		os.Exit(1)
	}
	//Print spaces
	spaces := GetSpaces()
	fmt.Println("Your Spaces:")
	for i, space := range spaces {
		fmt.Printf("[%d] %s\n", i, space.FriendlyName)
	}

	//Get Space Index
	fmt.Print("What Space would you like to delete? ")
	spaceIndexString, err := reader.ReadString('\n')
	spaceIndexString = strings.TrimSpace(spaceIndexString)
	spaceIndex, err := strconv.Atoi(spaceIndexString)
	if err != nil {
		fmt.Println("Invalid Choice: " + err.Error())
		os.Exit(1)
	}
	if spaceIndex < 0 || spaceIndex >= len(spaces) {
		fmt.Println("Invalid Choice: Space does not exist")
		os.Exit(1)
	}

	spaceToDelete := spaces[spaceIndex]
	//Confirm delete
	fmt.Printf("Are you sure you want to delete %s(%d) (yes,no)? ", spaceToDelete.FriendlyName, spaceToDelete.ID)
	confirmChoiceString, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Invalid Choice: " + err.Error())
		os.Exit(1)
	}
	if strings.TrimSpace(confirmChoiceString) == "yes" {
		fmt.Printf("Deleting %s (This may take a minute or two.)\n", spaceToDelete.FriendlyName)
		response, err := SendSpaceDeleteRequest(session, strconv.Itoa(int(spaceToDelete.ID))) //Yay for gross type conversions
		if err != nil {
			fmt.Println("Error sending request: " + err.Error())
			os.Exit(1)
		}
		if response.StatusCode != 200 {
			fmt.Println("Error Deleting Space!")
		} else {
			fmt.Println("Space Deleted.")
		}
	} else {
		fmt.Println("Not Deleting Space.")
	}
}

func SendSpaceDeleteRequest(session SessionRecord, spaceid string) (*http.Response, error) {
	//Build URL
	url := "https://" + session.OrchestratorHostname + "/api/v1/space/" + spaceid
	//Get HTTP Client
	hClient := GetHttpClient(session.IgnoreSSLErrors)
	//Create Request
	r, _ := http.NewRequest("DELETE", url, nil)
	r.Header.Add("X-Auth-Token", session.SessionToken)
	//Send the data and get the response
	resp, err := hClient.Do(r)
	if err != nil {
		fmt.Errorf("Error: %s\n", err.Error())
		return nil, err
	}

	return resp, nil
}

func init() {
	RootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
