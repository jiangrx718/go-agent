package viper

import (
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Init(configPath, envPath string) error {
	envPaths := make([]string, 0)
	if envPath != "" {
		envPaths = append(envPaths, envPath)
	}
	if err := godotenv.Load(envPaths...); err != nil && !os.IsNotExist(err) {
		return err
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config changed, reloading...")
	})

	return nil
}
