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
	"github.com/twa16/go-cas/client"
	"os"
	"bufio"
	"strings"
	"strconv"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to a Userspace cluster",
	Long: `Login to a Userspace cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		//Let's the hostname and orc info
		fmt.Print("Enter Orchestrator Address: ")
		orcHostname, _ := reader.ReadString('\n')
		orcHostname = strings.TrimSpace(orcHostname)
		orcInfo, err := GetOrchestratorInformation(orcHostname)
		if err != nil {
			fmt.Println("Error: "+err.Error())
			os.Exit(1)
		}
		var ignoreSSL bool
		for true {
			fmt.Print("Ignore SSL Errors(true/false): ")
			ignoreSSLString, err := reader.ReadString('\n')
			ignoreSSLString = strings.TrimSpace(ignoreSSLString)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				os.Exit(1)
			}
			ignoreSSL, err = strconv.ParseBool(ignoreSSLString)
			if err == nil {
				break
			}
		}

		//Make sure the orchestrator is allowing logins
		if !orcInfo.AllowsLocalLogin && !orcInfo.SupportsCAS {
			fmt.Errorf("Error: %s\n","Orchestrator not allowing logins")
			os.Exit(1)
		}

		useCAS := true
		//If there are multiple auth methods choose one
		if orcInfo.AllowsLocalLogin && orcInfo.SupportsCAS {
			fmt.Print("What authentication method do you use wish to use(CAS/Local): ")
			useCASString, _ := reader.ReadString('\n')
			useCASString = strings.TrimSpace(useCASString)
			useCASString = strings.ToLower(useCASString)
			if useCASString != "cas" {
				useCAS = false
			}
		}
		//If using CAS, go through that flow
		if useCAS {
			config := gocas.CASServerConfig{}
			config.ServerHostname = orcInfo.CASURL
			config.IgnoreSSLErrors = false

			ticket := config.StartLocalAuthenticationProcess()
			session, err := SubmitCASTicket(ticket)
			if err != nil {
				fmt.Println(err)
			}
			err = SaveSession(*session, orcHostname, ignoreSSL)
			if err != nil {
				fmt.Println("Failed to save session: "+ err.Error())
				panic("Could not save session. Send help.")
			}
			fmt.Println("Sucessfully logged in and saved session.")

		} else {
			fmt.Errorf("Error: %s\n","Local Login not yet Supported!")
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
