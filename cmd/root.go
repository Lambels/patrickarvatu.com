package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "backend-cli",
	Short: "Cli to start up patrickarvatu.com backend server",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./patrickarvatu.toml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// configure viper.
	viper.SetConfigName("patrickarvatu")
	viper.SetConfigType("toml")

	if cfgFile != "" {
		// use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// search config in home directory with name ".main" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".main")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// if a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
