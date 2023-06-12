package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"github.com/zembrodt/showtunes-api"
	"github.com/zembrodt/showtunes-api/controller"
	"github.com/zembrodt/showtunes-api/util/global"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var osEnvConfigs = []string{
	global.ServerAddressKey,
	global.ServerPortKey,
	global.ClientIdKey,
	global.ClientSecretKey,
	global.OriginKey,
	global.MaxAgeKey,
	global.ValidDomainsKey,
}

func main() {
	// CLI arguments
	var displayVersion bool
	flag.BoolVar(&displayVersion, "version", false, "Check version for this Spotify Auth API build")
	flag.Parse()

	if displayVersion {
		fmt.Printf("%s v%s\n", showtunes.Name, showtunes.Version)
		os.Exit(0)
	}

	setConfigurations()

	serverAddress := viper.GetString(global.ServerAddressKey)
	serverPort := viper.GetInt(global.ServerPortKey)

	server := controller.New(viper.GetString(global.ClientIdKey), viper.GetString(global.ClientSecretKey))
	server.Start(serverAddress, serverPort)
}

func setConfigurations() {
	// Set config defaults
	viper.SetDefault(global.ServerAddressKey, "localhost")
	viper.SetDefault(global.ServerPortKey, 8000)
	viper.SetDefault(global.ClientIdKey, "")
	viper.SetDefault(global.ClientSecretKey, "")
	viper.SetDefault(global.OriginKey, "*")
	viper.SetDefault(global.MaxAgeKey, "86400")
	viper.SetDefault(global.ValidDomainsKey, "i.scdn.co")

	// Get config from config.yaml
	viper.SetConfigName(global.ConfigFileName)
	viper.SetConfigType(global.ConfigFileExtension)
	_, b, _, _ := runtime.Caller(0)
	projectPath, _ := filepath.Split(filepath.Dir(b))
	viper.AddConfigPath(filepath.Join(projectPath, global.ConfigFilePath))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, continue with defaults
			log.Printf("%s.%s not found, using defaults\n",
				filepath.Join(global.ConfigFilePath, global.ConfigFileName),
				global.ConfigFileExtension,
			)
		} else {
			// Config file found, but error produced
			panic(fmt.Errorf("Fatal error reading config file: %s\n", err))
		}
	}

	// Check for environment variables
	for _, key := range osEnvConfigs {
		val, lookup := os.LookupEnv(global.EnvPrefix + "_" + key)
		if lookup {
			log.Printf("Adding env variable for %s\n", key)
			viper.Set(key, val)
		}
	}
}
