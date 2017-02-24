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
	"net/http"
	"bytes"
	"bufio"
	"log"
	"github.com/twa16/userspace/daemon"
	"encoding/json"
	"strings"
	"os"
	"strconv"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Space Request
		spaceRequest := userspaced.Space{}

		//Create reader for stdin
		reader := bufio.NewReader(os.Stdin)
		//Get the saved session
		session, err := GetSavedSession()
		if err != nil {
			fmt.Errorf("Error: %s\n", err.Error())
			return
		}
		//Display Possible Images to User
		images := GetImages()
		fmt.Println("Available Images:")
		for _, image := range images {
			if image.Active {
				fmt.Printf("%d %s\n", image.ID, image.Name)
			}
		}
		fmt.Println("Please choose an image. Run the images subcommand to see more details about available images.")
		/*
		 * Get image id
		 */
		//Get user input
		fmt.Print("Image ID: ")
		imageIDString, err := reader.ReadString('\n')
		imageIDString = strings.TrimSpace(imageIDString)
		if err != nil {
			fmt.Println("Failed to get input: "+err.Error())
			os.Exit(1)
		}
		//convert to int
		imageID, err := strconv.Atoi(imageIDString)
		if err != nil {
			fmt.Println("Invalid Image ID: Bad format")
			os.Exit(1)
		}
		//check if the id is that of a valid image
		isValidImage := false
		for _, image := range images {
			if image.ID == uint(imageID) {
				isValidImage = true
			}
		}
		//Report if the image is not valid
		if !isValidImage {
			fmt.Println("Invalid Image ID: Image does not exist")
			os.Exit(1)
		}
		//Save the image id
		spaceRequest.ImageID = uint(imageID)

		/*
		 * Get Space Name
		 */
		fmt.Print("Space Name: ")
		spaceName, _ := reader.ReadString('\n')
		spaceName = strings.TrimSpace(spaceName)
		spaceRequest.FriendlyName = spaceName

		/*
		 * Choose SSH Key
		 */
		//TODO: USE SSH KEYS

		/*
		 * Send Request
		 */
		createRespReader, err := SendSpaceCreateRequest(session, spaceRequest)
		if err != nil {
			fmt.Println("Error sending request: "+err.Error())
			os.Exit(1)
		}
		//Print response
		for {
			line, _ := createRespReader.ReadBytes('\n')
			if string(line) != ""{
				log.Println("\n"+string(line))
			}
			if strings.HasPrefix(string(line), "Error") ||
				strings.HasPrefix(string(line), "Creation Complete") ||
				strings.HasPrefix(string(line), "Session Expired") {
				os.Exit(1)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}

func SendSpaceCreateRequest(session SessionRecord, spaceRequest userspaced.Space) (*bufio.Reader, error){
	//Build URL
	url := "https://"+session.OrchestratorHostname+"/api/v1/spaces"
	//Get HTTP Client
	hClient := GetHttpClient(session.IgnoreSSLErrors)
	//JSONify Request
	jsonBytes, err := json.Marshal(&spaceRequest)
	if err != nil {
		fmt.Errorf("Error: %s\n", err.Error())
		return nil, err
	}
	//Create Request
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	r.Header.Add("X-Auth-Token", session.SessionToken)
	//Send the data and get the response
	resp, err := hClient.Do(r)
	if err != nil {
		fmt.Errorf("Error: %s\n", err.Error())
		return nil, err
	}
	//Get the body of the response as a string
	reader := bufio.NewReader(resp.Body)
	return reader, nil
}

func test() {
	// TODO: Work your own magic here
	fmt.Println("create called")



}
