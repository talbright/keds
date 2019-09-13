package config

import (
	"github.com/spf13/viper"
	"github.com/talbright/keds/utils/system"

	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

const (
	Name = "keds"
)

var ConfigPath = []string{".", fmt.Sprintf("$HOME/.%s", Name)}

func GetConfigPathKey() string {
	return fmt.Sprintf("%s_CONFIG_FILE_PATH", strings.ToUpper(Name))
}

func GetConfigPath() string {
	return os.Getenv(GetConfigPathKey())
}

func InitConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(Name)
		for _, p := range ConfigPath {
			viper.AddConfigPath(p)
		}
	}
	viper.SetEnvPrefix(Name)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("using config file: %s", viper.ConfigFileUsed())
		key := GetConfigPathKey()
		val := viper.ConfigFileUsed()
		os.Setenv(key, val)
	} else {
		if reflect.TypeOf(viper.ConfigFileNotFoundError{}) == reflect.TypeOf(err) {
			log.Printf("warning: no config file found in path (%s)\n", strings.Join(ConfigPath, ":"))
		} else {
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
	}
	expand()
}

func expand() {
	expandPluginPath()
}

func expandPluginPath() {
	pp := viper.GetStringSlice("plugin_path")
	if len(pp) > 0 {
		expanded := make([]string, 0)
		for _, v := range pp {
			expanded = append(expanded, system.AbsPathify(v))
		}
		viper.Set("plugin_path", expanded)
	}
}
