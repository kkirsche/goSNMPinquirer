// Copyright © 2016 Kevin Kirsche <kevin.kirsche@verizon.com> <kev.kirsche@gmail.com>
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
	"log"
	"os"

	"github.com/kkirsche/gosnmp"
	"github.com/spf13/cobra"
)

// walkCmd represents the walk command
var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "Walk is used to iterate through an OID tree",
	Long: `The walk command  performs a whole series of getnexts automatically
for you, and stops when it returns results that are no longer inside the range
of the OID which you originally specified. If you wanted to get all of the
information stored on a machine in the system MIB group, for instance, you
could use this command to do so. This will be slower than a bulk walk though.`,
	Run: func(cmd *cobra.Command, args []string) {
		if oid == "" {
			log.Fatal("Please provide an OID to retrieve")
			os.Exit(1)
		}

		snmp, err := gosnmp.Connect(remoteIP, community, gosnmp.Version2c, 50)
		if err != nil {
			log.Fatal(err.Error())
		}

		resultCh := make(chan gosnmp.SnmpPDU)

		go func(snmp *gosnmp.Conn, resultCh chan gosnmp.SnmpPDU) {
			snmp.StreamWalk(oid, resultCh)
		}(snmp, resultCh)

		for item := range resultCh {
			fmt.Println("Name:", item.Name, "Type:", item.Type, "Value:", item.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(walkCmd)
	walkCmd.PersistentFlags().StringVarP(&oid, "oid", "o", "", "The OID to retrieve")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// walkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// walkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
