// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"log"

	"k8s-webhook-admission-controller/pkg/certificate"

	"github.com/spf13/cobra"
)

var (
	writeDir string
	ip       string
	dns       string
	genClientCert bool
	genCaCert bool
	genServerCert bool
)

// initcertCmd represents the initcert command
var initcertCmd = &cobra.Command{
	Use:   "initcert",
	Short: "Create CA cert, server.cert, server.key, client.cert, client.key",
	Run: func(cmd *cobra.Command, args []string) {
		err := certificate.GenerateCertificate(writeDir, ip,dns,genCaCert,genServerCert,genClientCert)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(initcertCmd)
	initcertCmd.PersistentFlags().StringVar(&writeDir, "dir", "/home/ac/go/src/k8s-webhook-admission-controller", "directory to write cert files")
	initcertCmd.PersistentFlags().StringVar(&ip, "ip", "", "ip of server, required for server certificate")
	initcertCmd.PersistentFlags().StringVar(&dns, "dns", "", "dns of server, required for server certificate")
	initcertCmd.PersistentFlags().BoolVar(&genClientCert, "client", false, "generate client cert")
	initcertCmd.PersistentFlags().BoolVar(&genCaCert, "ca", false, "generate ca cert")
	initcertCmd.PersistentFlags().BoolVar(&genServerCert, "server", false, "generate server cert")
}
