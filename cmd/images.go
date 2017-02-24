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
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"os"
	"github.com/olekukonko/tablewriter"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Active", "ID", "Name", "Description"})
		for _, image := range GetImages() {
			line := []string {
				strconv.FormatBool(image.Active),
				strconv.Itoa(int(image.ID)),
				image.Name,
				image.Description,
			}
			table.Append(line)
		}
		table.Render()
	},
}

func init() {
	RootCmd.AddCommand(imagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func GetImages() []userspaced.SpaceImage {
	var images []userspaced.SpaceImage
	session, err := GetSavedSession()
	if err != nil {
		fmt.Println("Session Error: "+err.Error())
		return images
	}
	url := "https://"+session.OrchestratorHostname+"/api/v1/images"
	hClient := GetHttpClient(session.IgnoreSSLErrors)
	r, _ := http.NewRequest("GET", url, nil)
	r.Header.Add("X-Auth-Token", session.SessionToken)
	resp, err := hClient.Do(r)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	//TODO: Handle error conditions

	json.Unmarshal(buf.Bytes(), &images)
	return images
}