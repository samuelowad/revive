package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type configStruct struct {
	Command                string   `yaml:"command" json:"command"`
	IgnoreDirectories      []string `yaml:"ignoreDirectories" json:"ignoreDirectories"`
	MonitorFileExt         []string `yaml:"monitorFileExt" json:"monitorFileExt"`
	RestartDelaySeconds    int      `yaml:"restartDelaySeconds" json:"restartDelaySeconds"`
	IgnoreFileNameEndsWith []string `yaml:"ignoreFileNameEndsWith" json:"ignoreFileNameEndsWith"`
	VerbosityLevel         int      `yaml:"verbosityLevel" json:"verbosityLevel"`
}

var Config *configStruct

func parseYAMLConfig(configFile string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &Config); err != nil {
		return err
	}

	return nil
}

func parseJSONConfig(configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &Config); err != nil {
		return err
	}

	return nil
}

// Find config file
func findConfigFile() (string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.Name() == "revive.yaml" || file.Name() == "revive.json" {
			return file.Name(), nil
		}
	}

	log.Fatal("Config file 'revive.yaml' or 'revive.json' not found in the current directory")
	return "", nil
}

func ReadConfig() {
	configFile, err := findConfigFile()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Config file found:", configFile)

	if strings.HasSuffix(configFile, ".yaml") {
		parseYAMLConfig(configFile)
	} else if strings.HasSuffix(configFile, ".json") {
		parseJSONConfig(configFile)
	} else {
		log.Fatal("Unsupported config file format")
		return
	}

	if err != nil {
		log.Fatal("Error parsing config file:", err)
		return
	}

	fmt.Println("Config data:", *Config)
}
