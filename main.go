package main

import (
	// "github.com/spf13/viper"
	"github.com/talbright/keds/cmd"
	// . "github.com/talbright/keds/server"
	_ "github.com/talbright/keds/utils/config"
)

func main() {
	cmd.Execute()
	cmd.RootCmd.DebugFlags()
	// if err := cmd.RootCmd.Execute(); err != nil {
	// 	log.Println("error root cmd")
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// log.Printf("config: %v", viper.AllSettings())
	//gRPC := NewKedsRPCServer()
	//gRPC.Start()
}
