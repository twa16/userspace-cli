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
	"github.com/k0kubun/pp"
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
		// TODO: Work your own magic here
		fmt.Println("list called")

		pp.Println(GetSpaces())
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
	var spaces []userspaced.Space
	session, err := GetSavedSession()
	if err != nil {
		fmt.Println("Session Error: "+err.Error())
		return spaces
	}
	url := "https://"+session.OrchestratorHostname+"/api/v1/spaces"
	hClient := GetHttpClient(true)
	r, _ := http.NewRequest("GET", url, nil)
	r.Header.Add("X-Auth-Token", session.SessionToken)
	resp, err := hClient.Do(r)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	//TODO: Handle error conditions

	json.Unmarshal(buf.Bytes(), &spaces)
	return spaces
}
