package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	musicapi "github.com/zembrodt/music-display-api"
	"github.com/zembrodt/music-display-api/controller"
	"github.com/zembrodt/music-display-api/util/global"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var osEnvConfigs = []string{
	global.ServerAddress,
	global.ServerPort,
	global.PublicAddress,
	global.ClientIdKey,
	global.ClientSecretKey,
	global.OriginKey,
	global.MaxAgeKey,
}

func main() {
	// CLI arguments
	var displayVersion bool
	flag.BoolVar(&displayVersion, "version", false, "Check version for this Spotify Auth API build")
	flag.Parse()

	if displayVersion {
		fmt.Printf("%s v%s\n",musicapi.Name, musicapi.Version)
		os.Exit(0)
	}

	setConfigurations()

	serverAddress := viper.GetString(global.ServerAddress)
	serverPort := viper.GetInt(global.ServerPort)

	//repo := repository.New(db)
	//svc := service.New(repo, expiryTime)

	server := controller.New(viper.GetString(global.ClientIdKey), viper.GetString(global.ClientSecretKey))
	server.Start(serverAddress, serverPort)
}

func setConfigurations() {
	// Set config defaults
	viper.SetDefault(global.ServerAddress, "localhost")
	viper.SetDefault(global.ServerPort, 8000)
	viper.SetDefault(global.ClientIdKey, "")
	viper.SetDefault(global.ClientSecretKey, "")
	viper.SetDefault(global.MaxAgeKey, "86400")

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
		val, lookup := os.LookupEnv(key)
		if lookup {
			log.Printf("Adding env variable for %s\n", key)
			viper.Set(key, val)
		}
	}
}
