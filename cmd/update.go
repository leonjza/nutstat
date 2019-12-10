package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/leonjza/nutstat/simpleflux"

	nut "github.com/robbiet480/go.nut"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates UPS statistics",
	Long: `Connects to an InfluxDB instance and writes
updated UPS statistics at an interval.`,
	Run: func(cmd *cobra.Command, args []string) {

		// connect nut
		client, err := nut.Connect(viper.GetString("nuthost"))
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.Authenticate(viper.GetString("nutusername"),
			viper.GetString("nutpassword"))
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Connected to NUT @ %s", viper.GetString("nuthost"))

		// setup influxDB
		c := simpleflux.NewSimpleInfluxDB()

		log.Printf("Making sure Influx is reachable")
		for {
			status, err := c.Ping()
			if status == true {
				break
			}
			log.Printf("InfluxDB is not available: %s. Sleeping 30s", err)
			time.Sleep(time.Second * 30)
		}

		// search for the UPS we should poll
		upsList, err := client.GetUPSList()

		var targetUps *nut.UPS
		for _, u := range upsList {

			if u.Name == viper.GetString("nutupsname") {
				targetUps = &u
			}
		}

		if targetUps != nil {
			log.Printf(`Found UPS: %s`, viper.GetString("nutupsname"))
		} else {
			log.Fatalf(`Could not find UPS named: %s`, viper.GetString("nutupsname"))
		}

		log.Printf(`Processing updates every %d seconds...`, viper.GetInt32("interval"))

		for {
			time.Sleep(time.Second * time.Duration(viper.GetInt32("interval")))

			// var statusLine string
			values, _ := targetUps.GetVariables()
			statusLine := `ups,ups_name=` + viper.GetString("nutupsname") + ` `
			sanityCheck := statusLine

			for _, v := range values {
				switch v.Name {
				case "battery.charge":
					statusLine += fmt.Sprintf(`%s=%d,`, v.Name, v.Value.(int64))
				case "battery.runtime":
					statusLine += fmt.Sprintf(`%s=%d,`, v.Name, v.Value.(int64))
				case "input.voltage":
					statusLine += fmt.Sprintf(`%s=%.2f,`, v.Name, v.Value.(float64))
				case "output.frequency":
					statusLine += fmt.Sprintf(`%s=%.2f,`, v.Name, v.Value.(float64))
				case "output.frequency.nominal":
					statusLine += fmt.Sprintf(`%s=%d,`, v.Name, v.Value.(int64))
				case "output.voltage":
					statusLine += fmt.Sprintf(`%s=%.2f,`, v.Name, v.Value.(float64))
				case "ups.load":
					statusLine += fmt.Sprintf(`%s=%d,`, v.Name, v.Value.(int64))
				case "ups.status":
					statusLine += fmt.Sprintf(`%s="%s",`, v.Name, v.Value.(string))
				}
			}

			statusLine = strings.TrimSuffix(statusLine, ",")
			if statusLine == sanityCheck {
				log.Fatalf(`Could cound read information from NUT. Exiting...`)
			}
			log.Printf(`Posting line: %s`, statusLine)

			if err = c.Write(statusLine); err != nil {
				log.Printf(`Error posting to InfluxDB: %s`, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().Int32P("interval", "i", 30, "Stats update interval in seconds")
	viper.BindPFlag("interval", updateCmd.PersistentFlags().Lookup("interval"))
}
