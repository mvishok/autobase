package loader

import (
	"encoding/json"
	"os"
	"statix/pkg/core/memory"
	"statix/pkg/log"
)

func readJson(filename string, v interface{}) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(fileData, v)
}

func LoadConfig(path string) map[string]interface{} {
	var config interface{}
	err := readJson(path, &config)
	if err != nil {
		// Handle the error
		log.Error("Failed to read config file: " + err.Error())
		// Return an empty map
		return make(map[string]interface{})
	}

	configMap, ok := config.(map[string]interface{})
	if !ok {
		log.Error("Invalid config format")
		return make(map[string]interface{})
	}

	//check if config["dir"] is set
	if configMap["dir"] == nil {
		log.Error("No dir set in config")
		os.Exit(1)
	}

	disableLogging, ok := configMap["disableLogging"].(bool)
	if ok {
		if disableLogging {
			log.Disable()
		} else {
			log.Enable()
		}
	} else {
		log.Enable()
	}

	// Set the default port if it's not set in the config
	if configMap["port"] == nil {
		configMap["port"] = 8081
		log.Warning("Port not set in config, defaulting to 8081")
	}

	memory.Append("test", 12) //REMOVE

	return configMap
}
