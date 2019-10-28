package cmd

import (
	"log"

	"github.com/leonjza/nutstat/simpleflux"
	nut "github.com/robbiet480/go.nut"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if your settings seem to work",
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
		c.Ping()

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

	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
