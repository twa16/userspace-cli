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
		// TODO: Work your own magic here
		fmt.Println("create called")

		//Get the saved session
		session, err := GetSavedSession()
		if err != nil {
			fmt.Errorf("Error: %s\n", err.Error())
			return
		}

		//Build URL
		url := "https://"+session.OrchestratorHostname+"/api/v1/spaces"
		//Get HTTP Client
		hClient := GetHttpClient(true)
		//Create Space Object
		spaceRequest := userspaced.Space{}
		spaceRequest.FriendlyName = "Test Space"
		spaceRequest.SSHKeyID = 0
		spaceRequest.ImageID = 1
		//JSONify Request
		jsonBytes, err := json.Marshal(&spaceRequest)
		if err != nil {
			fmt.Errorf("Error: %s\n", err.Error())
			return
		}
		//Create Request
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
		r.Header.Add("X-Auth-Token", session.SessionToken)
		//Send the data and get the response
		resp, err := hClient.Do(r)
		if err != nil {
			fmt.Errorf("Error: %s\n", err.Error())
			return
		}
		//Get the body of the response as a string
		reader := bufio.NewReader(resp.Body)
		for {
			line, _ := reader.ReadBytes('\n')
			if string(line) != ""{
				log.Println("\n"+string(line))
			}
			if strings.HasPrefix(string(line), "Error") || strings.HasPrefix(string(line), "Creation Complete") {
				os.Exit(1)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
