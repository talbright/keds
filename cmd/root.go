package cmd

import (
	"fmt"
	"log"
	"os"

	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/talbright/keds/server"
	"github.com/talbright/keds/utils/config"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "keds",
	Short: "A prototype for a generic CLI plugin framework based on gRPC.",
	Long:  "See http://github.com/talbright/keds/README.md",
	PreRun: func(cmd *cobra.Command, args []string) {
		log.Printf("Cobra.PreRun")
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Cobra.Run")
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.keds.yaml)")

	// cobra doesn't invoke these callbacks until *after* checking for the help flag, so
	// we can't use this...
	// see https://github.com/spf13/cobra/blob/6ed17b5128e8932c9ecd4c3970e8ea5e60a418ac/command.go#L590
	// cobra.OnInitialize(onInitialize)
	onInitialize()
	log.Printf("config: %v", viper.AllSettings())
}

func onInitialize() {
	initializeConfig()
	initializeRuntime()
}

func initializeConfig() {
	config.InitConfig(cfgFile)
}

func initializeRuntime() {
	go func() {
		gRPC := server.NewKedsRPCServer()
		gRPC.Cobra = server.NewCobra(RootCmd)
		gRPC.Start()
	}()
	//TODO this feels hacky...need a more reliable way to determine that the plugins have loaded
	//and the server has started
	time.Sleep(3 * time.Second)
}
