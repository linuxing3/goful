package config

import (
	"fmt"
	"os"

	_ "embed"

	"github.com/spf13/viper"
)

//go:embed config.yaml
var ConfigTemplate string

func InitConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name

	path := fmt.Sprintf("%s/.goful", os.Getenv("HOME")) // call multiple times to add many search paths
	searchPath := []string{
		path,
		".",
		"/etc/goful",
		"~/.goful",
	}
	for _, path := range searchPath {
		viper.AddConfigPath(path) // call multiple times to add many search paths
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		MakeDefaultConfig(ConfigTemplate)
	}
}

func MakeDefaultConfig(text string) {
	path := fmt.Sprintf("%s/.goful/config.yaml", os.Getenv("HOME")) // call multiple times to add many search paths
	c, err := os.ReadFile(path)
	if string(c) == "" || err != nil {
		fmt.Println("Config file not exists or is empty, fix for you...")
		if err := os.WriteFile(path, []byte(text), 644); err != nil {
			fmt.Println("Please restart your goful!")
		}
	}
}
