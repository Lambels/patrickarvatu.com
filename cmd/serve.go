package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/Lambels/patrickarvatu.com/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the application",
	RunE:  RunServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().BoolP("debug", "d", false, "starts a prom-http on port :8000 /metrics.")
}

func RunServe(cmd *cobra.Command, _ []string) error {
	var cfg pa.Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	_, cleanUp, err := initializeServer(&cfg)
	if err != nil {
		return err
	}
	defer cleanUp()

	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		go http.RunDebugServer()
		log.Println("[INFO] Started debug server on :8000 /metrics.")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	sig := <-c
	log.Println("[INFO] Got signal", sig.String(), ", exiting.")
	return nil
}
