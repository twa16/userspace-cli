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
	"github.com/twa16/userspace/daemon"
	"net/http"
	"bytes"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Image", "Name", "State", "SSH Host", "SSH Port"})
		for _, space := range GetSpaces() {
			var sshHost string
			var sshPort uint16
			for _, mapping := range space.PortLinks {
				if mapping.SpacePort == 22 {
					sshPort = mapping.ExternalPort
					sshHost = mapping.ExternalAddress
				}
			}
			line := []string {
				strconv.Itoa(int(space.ID)),
				getImageNameByID(space.ID),
				space.FriendlyName,
				space.SpaceState,
				sshHost,
				strconv.FormatUint(uint64(sshPort), 10),
			}
			table.Append(line)
		}
		table.Render()
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func GetSpaces() []userspaced.Space {
	spaces := []userspaced.Space{}
	session, err := GetSavedSession()
	if err != nil {
		fmt.Println("Session Error: "+err.Error())
		return spaces
	}
	url := "https://"+session.OrchestratorHostname+"/api/v1/spaces"
	hClient := GetHttpClient(true)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Errorf("Error:  %s\n", err.Error())
		return spaces
	}
	r.Header.Add("X-Auth-Token", session.SessionToken)
	resp, err := hClient.Do(r)
	if err != nil {
		fmt.Errorf("Error Getting Spaces: %s\n", err.Error())
		return spaces
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	if buf.String() == "Session Expired" {
		fmt.Println("Session is Expired")
		os.Exit(2)
	}
	//TODO: Handle error conditions

	json.Unmarshal(buf.Bytes(), &spaces)
	return spaces
}

func getImageNameByID(id uint) string{
	for _, image := range GetImages() {
		if image.ID == id {
			return image.Name
		}
	}
	return "ERROR"
}