/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	appSrv "github.com/dengzitong/probe-host/cmd/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	addr          string
	port          string
	probe_timeout int
)

func init() {
	pflag.StringVar(&addr, "addr", "127.0.0.1", "--addr 127.0.0.1")
	pflag.StringVar(&port, "port", "9999", "--port 9999")
	pflag.IntVar(&probe_timeout, "probe-timeout", 1, "--probe-timeout 1")
	pflag.Parse()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "probehost -addr 127.0.0.1 -port 9999 -probe-timeout 1",
	Short: "A probe host http service",
	Long:  `Implementing an http request responseprocess needs to detect whether the ip address port can be reached.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return appSrv.RunServer(addr, port, probe_timeout)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
