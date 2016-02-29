// Copyright © 2016 Kevin Kirsche <kev.kirsche@gmail.com>
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
	"time"

	"github.com/kkirsche/gosnmp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dayCmd represents the day command
var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "execute the daily cron job",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		snmp, err := gosnmp.Connect(viper.GetString("ip"), viper.GetString("community"), gosnmp.Version2c, 50)
		if err != nil {
			log.Fatal(err.Error())
		}

		var hostname string

		getValues := viper.GetStringSlice("cron.day.get")
		for _, oid := range getValues {
			pdu, err := snmp.Get(oid)
			if err != nil {
				log.Fatal(err.Error())
			}

			fmt.Print(time.Now().UTC().Format(time.RFC3339Nano))
			fmt.Print(",", snmp.Target)
			if hostname != "" {
				fmt.Print(",", hostname)
			}

			for _, variable := range pdu.Variables {
				fmt.Print(",", variable.Name)
				fmt.Print(",", variable.Value)
			}
			fmt.Print("\n")

			if oid == ".1.3.6.1.2.1.1.5.0" {
				hostname = pdu.Variables[0].Value.(string)
			}
		}

		getBulkValues := viper.GetStringSlice("cron.day.getbulk")
		for _, oid := range getBulkValues {
			pdu, err := snmp.GetBulk(0, 100, oid)
			if err != nil {
				log.Fatal(err.Error())
			}

			fmt.Print(time.Now().UTC().Format(time.RFC3339Nano))
			fmt.Print(",", snmp.Target)
			if hostname != "" {
				fmt.Print(",", hostname)
			}

			for _, pdu := range pdu.Variables {
				fmt.Print(",", pdu.Name)
				fmt.Print(",", pdu.Value)
			}
			fmt.Print("\n")
		}

		bulkwalkValues := viper.GetStringSlice("cron.day.bulkwalk")
		results := make(map[string][]gosnmp.SnmpPDU)
		lastOID := ""
		for _, oid := range bulkwalkValues {
			pdus, err := snmp.BulkWalk(100, oid)
			if err != nil {
				log.Println("Error: ", err.Error())
			}

			results[oid] = pdus
			lastOID = oid
		}

		lengthOfValues := len(results[lastOID]) - 1
		for i := 0; i < lengthOfValues; i++ {
			fmt.Print(time.Now().UTC().Format(time.RFC3339Nano))
			fmt.Print(",", snmp.Target)
			if hostname != "" {
				fmt.Print(",", hostname)
			}
			for _, oid := range bulkwalkValues {
				oidLength := len(results[oid]) - 1
				if oidLength >= i {
					pdu := results[oid][i]
					fmt.Print(",", pdu.Name)
					fmt.Print(",", pdu.Value)
				}
			}
			fmt.Print("\n")
		}
	},
}

func init() {
	cronCmd.AddCommand(dayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
