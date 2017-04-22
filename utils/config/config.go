package config

import (
	"github.com/spf13/viper"

	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

const (
	NAME = "keds"
)

var CONFIG_PATH = []string{".", fmt.Sprintf("$HOME/.%s", NAME)}

func init() {
	viper.SetEnvPrefix(NAME)
	viper.AutomaticEnv()
	viper.SetConfigName(NAME)
	for _, p := range CONFIG_PATH {
		viper.AddConfigPath(p)
	}
	if err := viper.ReadInConfig(); err != nil {
		if reflect.TypeOf(viper.ConfigFileNotFoundError{}) == reflect.TypeOf(err) {
			log.Printf("warning: no config file found in path (%s)\n", strings.Join(CONFIG_PATH, ":"))
		} else {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	if err := expand(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func expand() (err error) {
	expandPluginPath()
	return nil
}

func expandPluginPath() {
	pp := viper.GetStringSlice("plugin_path")
	if len(pp) > 0 {
		expanded := make([]string, 0)
		for _, v := range pp {
			expanded = append(expanded, AbsPathify(v))
		}
		viper.Set("plugin_path", expanded)
	}
}

func AbsPathify(inPath string) string {

	if strings.HasPrefix(inPath, "$HOME") {
		inPath = UserHomeDir() + inPath[5:]
	}

	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}

	return ""
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
