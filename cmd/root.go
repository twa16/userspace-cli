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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"crypto/tls"
	"net/http"
	"github.com/twa16/userspace/daemon"
	"encoding/json"
	"bytes"
	"github.com/twa16/go-auth"
	"github.com/pkg/errors"
	"strconv"
	"io/ioutil"
	"github.com/mitchellh/go-homedir"
)

type SessionRecord struct {
	UserID uint
	SessionToken string
	OrchestratorHostname string
	IgnoreSSLErrors bool
}

var cfgFile string
var OrchestratorInfo *userspaced.OrchestratorInfo
var OrchestratorURL string
// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "userspace-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
// Uncomment the following line if your bare application
// has an action associated with it:
//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.userspace-cli.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".userspace-cli") // name of config file (without extension)
	viper.AddConfigPath("$HOME")          // adding home directory as first search path
	viper.AutomaticEnv()                  // read in environment variables that match
	homeDirectory, _ := homedir.Dir()
	viper.SetDefault("HomeDirectory", homeDirectory)
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func GetHttpClient(ignoreSSLErrors bool) *http.Client {
	//Create the client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: ignoreSSLErrors},
	}
	return &http.Client{Transport: tr}
}

func GetOrchestratorInformation(url string) (*userspaced.OrchestratorInfo, error) {
	hClient := GetHttpClient(true)
	resp, err := hClient.Get("https://"+url+"/orchestratorinfo")
	if err != nil {
		return nil, err
	}
	//Get the body of the response as a string
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	var orcInfo *userspaced.OrchestratorInfo
	err = json.Unmarshal(buf.Bytes(), &orcInfo)
	if err != nil {
		return nil, err
	}
	OrchestratorInfo = orcInfo
	OrchestratorURL = "https://"+url
	return orcInfo, nil
}

//SubmitCASTicket Attempts authentication against orchestrator
func SubmitCASTicket(ticket string) (*simpleauth.Session, error) {
	hClient := GetHttpClient(true)
	resp, err := hClient.Get(OrchestratorURL+"/caslogin?ticket="+ticket)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	//See if we logged in successfully
	if resp.StatusCode == 401 {
		return nil, errors.New("Invalid Login")
	}

	//Get the body of the response as a string
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	//Catch other errors
	if resp.StatusCode != 200 {
		return nil, errors.New("Error "+strconv.Itoa(resp.StatusCode)+": "+buf.String())
	}

	//Unmarshal our session
	var session simpleauth.Session
	err = json.Unmarshal(buf.Bytes(), &session)
	if err != nil {
		return nil, err
	}

	//Ok let's return our session
	return &session, nil
}

func SaveSession(sessionObj simpleauth.Session, hostname string, ignoreSSLErrors bool) error {
	sessionRecord := SessionRecord{
		SessionToken: sessionObj.AuthenticationToken,
		UserID: sessionObj.AuthUserID,
		OrchestratorHostname: hostname,
		IgnoreSSLErrors: ignoreSSLErrors,
	}
	sessionJSON, _ := json.Marshal(sessionRecord)
	err := ioutil.WriteFile(viper.GetString("HomeDirectory")+"/.userspace-session", sessionJSON, 0644)
	return err
}

func GetSavedSession() (SessionRecord, error) {
	var sessionRecord SessionRecord
	sessionBytes, err := ioutil.ReadFile(viper.GetString("HomeDirectory") + "/.userspace-session")
	if err != nil {
		return sessionRecord, err
	}
	err = json.Unmarshal(sessionBytes, &sessionRecord)
	return sessionRecord, err
}
