package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "nutstat",
	Short: "Push NUT statistics to InfluxDB",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nutstat.yaml)")

	// NUT config
	rootCmd.PersistentFlags().String("nuthost", "127.0.0.1", "The NUT host to connect to")
	viper.BindPFlag("nuthost", rootCmd.PersistentFlags().Lookup("nuthost"))
	rootCmd.PersistentFlags().String("nutusername", "", "The NUT username")
	viper.BindPFlag("nutusername", rootCmd.PersistentFlags().Lookup("nutusername"))
	rootCmd.PersistentFlags().String("nutpassword", "127.0.0.1", "The NUT password")
	viper.BindPFlag("nutpassword", rootCmd.PersistentFlags().Lookup("nutpassword"))
	rootCmd.PersistentFlags().String("nutupsname", "undefined", "UPS Name to use for stats")
	viper.BindPFlag("nutupsname", rootCmd.PersistentFlags().Lookup("nutupsname"))

	// InfluxDB config
	rootCmd.PersistentFlags().String("influxdbhost", "http://localhost:8086", "The InfluxDB host to connect to")
	viper.BindPFlag("influxdbhost", rootCmd.PersistentFlags().Lookup("influxdbhost"))
	rootCmd.PersistentFlags().String("influxdbdatabase", "ups", "The InfluxDB database to use")
	viper.BindPFlag("influxdbdatabase", rootCmd.PersistentFlags().Lookup("influxdbdatabase"))
	rootCmd.PersistentFlags().String("influxdbusername", "", "The InfluxDB username")
	viper.BindPFlag("influxdbusername", rootCmd.PersistentFlags().Lookup("influxdbusername"))
	rootCmd.PersistentFlags().String("influxdbpassword", "", "The InfluxDB password")
	viper.BindPFlag("influxdbpassword", rootCmd.PersistentFlags().Lookup("influxdbpassword"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".nutstat" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".nutstat")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
