// Copyright © 2018 Yale University
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/YaleUniversity/go-filelocker/pkg/filelocker"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, userID, apiKey, filelockerURL, clientTimeout string
var filelockerClient *filelocker.Client
var asJSON bool

// Version is the main version number
const Version = filelocker.Version

// VersionPrerelease is a prerelease marker
const VersionPrerelease = filelocker.VersionPrerelease

// Logger is a STDERR logger
var Logger = log.New(os.Stderr, "", 0)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "filelocker",
	Short: "Filelocker 2 client.",
	Long:  `A go cli for interacting with filelocker 2.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if filelockerURL == "" {
			return errors.New("filelocker URL is required")
		}

		t, err := time.ParseDuration(clientTimeout)
		if err != nil {
			return errors.New("cannot parse client timeout")
		}

		httpClient := &http.Client{
			Timeout: t * time.Second,
		}

		filelockerClient, err = filelocker.NewClient(userID, apiKey, filelockerURL, httpClient)
		return err
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "filelocker config file -- _not_ the control file (default is $HOME/.filelocker.yaml)")
	RootCmd.PersistentFlags().StringVarP(&userID, "login", "l", "", "The userid to use for connections to filelocker")
	RootCmd.PersistentFlags().StringVarP(&clientTimeout, "timeout", "t", "30s", "The filelocker http client timeout (seconds)")
	RootCmd.PersistentFlags().StringVarP(&apiKey, "key", "k", "", "The api key to use for connections to filelocker")
	RootCmd.PersistentFlags().StringVarP(&filelockerURL, "url", "u", "", "The base URL to use for connections to filelocker (ie. https://files.yale.edu")
	RootCmd.PersistentFlags().BoolVarP(&asJSON, "json", "j", false, "Format the response as JSON where applicable")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".filelocker")     // name of config file (without extension)
	viper.AddConfigPath(os.Getenv("HOME")) // adding home directory as first search path

	viper.SetEnvPrefix("filelocker") // prefix environment variables with FILELOCKER
	viper.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		Logger.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		Logger.Println(err)
		os.Exit(-1)
	}
}
