package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	_ "github.com/spf13/viper"
	"github.com/talbright/keds/utils/config"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "keds",
	Short: "A prototype for a generic and opinionated CLI plugin framework.",
	Long:  "See http://github.com/talbright/keds/README.md",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO start server, load plugins, etc
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
	cobra.OnInitialize(func() { config.InitConfig(cfgFile) })

	// Define persistent flags global for the application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.keds.yaml)")

	// Define local flags only run when this action is called directly.
	// RootCmd.Flags().BoolP("debug", "t", false, "Enable debug mode")
}
