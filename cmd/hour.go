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
	"log/syslog"
	"os"
	"strings"
	"time"

	"github.com/kkirsche/gosnmp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// hourCmd represents the hour command
var hourCmd = &cobra.Command{
	Use:   "hour",
	Short: "execute the cron job. For use once per hour",
	Long: `A job which should be executed once per hour. With crontab, this would
be executed using either the @daily string or the crontab definition "0 0	* * *".
This command loads from the config file (by default located at
$HOME/.inquirer.json) only and does not accept command line arguments except for
IP and Community String`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		saveMethod = strings.ToLower(viper.GetString("cron.save_via"))
		switch saveMethod {
		case "file":
			savePath := viper.GetString("cron.save_filepath")
			if savePath != "" {
				file, err = os.Create(savePath + viper.GetString("cron.save_filename") + time.Now().UTC().Format(time.RFC3339) + "-daily.csv")
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		case "syslog":
			syslogger, err = syslog.New(syslog.LOG_INFO, "Inquirer | Hour")
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		snmp, err := gosnmp.Connect(viper.GetString("ip"), viper.GetString("community"), gosnmp.Version2c, 30)
		if err != nil {
			syslogger.Err(err.Error())
			log.Fatal(err.Error())
		}

		var line string
		var hostname string
		getValues := viper.GetStringSlice("cron.hour.get")
		for _, oid := range getValues {
			pdu, err := snmp.Get(oid)
			if err != nil {
				syslogger.Err(err.Error())
				log.Fatal(err.Error())
			}

			line = time.Now().UTC().Format(time.RFC3339Nano) + "," + snmp.Target
			if hostname != "" {
				line += "," + hostname
			}

			for _, variable := range pdu.Variables {
				line += "," + variable.Name + "," + fmt.Sprintf("%v", variable.Value)
			}

			writeToOutputMethod(line + "\n")

			if oid == ".1.3.6.1.2.1.1.5.0" {
				hostname = pdu.Variables[0].Value.(string)
			}
		}

		getBulkValues := viper.GetStringSlice("cron.hour.getbulk")
		for _, oid := range getBulkValues {
			pdu, err := snmp.GetBulk(0, 100, oid)
			if err != nil {
				syslogger.Err(err.Error())
				log.Fatal(err.Error())
			}

			line = time.Now().UTC().Format(time.RFC3339Nano) + "," + snmp.Target
			if hostname != "" {
				line += "," + hostname
			}

			for _, pdu := range pdu.Variables {
				line += "," + pdu.Name + "," + fmt.Sprintf("%v", pdu.Value)
			}
			writeToOutputMethod(line + "\n")
		}

		bulkwalkValues := viper.GetStringSlice("cron.hour.bulkwalk")
		results := make(map[string][]gosnmp.SnmpPDU)
		for _, oid := range bulkwalkValues {
			pdus, err := snmp.BulkWalk(100, oid)
			if err != nil {
				syslogger.Err(err.Error())
				log.Println("Error: ", err.Error())
			}

			results[oid] = pdus
		}

		var lengthOfValues int
		for _, value := range results {
			currentOIDLength := len(value)
			if currentOIDLength > lengthOfValues {
				lengthOfValues = currentOIDLength
			}
		}

		for i := 0; i < lengthOfValues; i++ {
			line = time.Now().UTC().Format(time.RFC3339Nano) + "," + snmp.Target
			if hostname != "" {
				line += "," + hostname
			}
			for _, oid := range bulkwalkValues {
				oidLength := len(results[oid]) - 1
				if oidLength >= i {
					pdu := results[oid][i]
					line += "," + pdu.Name + "," + fmt.Sprintf("%v", pdu.Value)
				}
			}
			writeToOutputMethod(line + "\n")
		}
	},
}

func init() {
	cronCmd.AddCommand(hourCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hourCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hourCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
