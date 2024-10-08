package loader

import (
	"autobase/pkg/log"
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	"github.com/joho/godotenv"
)

func normalize(fpath string) string { //normalize the path
	return path.Clean(filepath.ToSlash(fpath))
}

func readJson(filename string, v interface{}) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(fileData, v)
}

func LoadConfig(filePath string) map[string]interface{} {
	log.Info("Loading config file: " + filePath)

	abs, absErr := filepath.Abs(filePath)
	if absErr != nil {
		log.Error("Failed to get absolute path: " + absErr.Error())
		os.Exit(1)
	}

	filePath = normalize(abs)

	var config interface{}
	err := readJson(filePath, &config)
	if err != nil {
		// Handle the error
		log.Error("Failed to read config file: " + err.Error())
		os.Exit(1)
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
	} else {
		configMap["dir"] = normalize(path.Join(path.Dir(filePath), configMap["dir"].(string))) //resolve the path

		if _, err := os.Stat(configMap["dir"].(string)); os.IsNotExist(err) {
			log.Error("Directory does not exist: " + configMap["dir"].(string))
			os.Exit(1)
		}
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

	//if env is set in config, load the env file
	if configMap["env"] != nil {
		envFilePath := normalize(path.Join(path.Dir(filePath), configMap["env"].(string)))
		if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
			log.Error("Env file does not exist: " + envFilePath)
			os.Exit(1)
		}

		err := godotenv.Load(envFilePath)
		if err != nil {
			log.Error("Error loading env file: " + err.Error())
			os.Exit(1)
		}
	}

	return configMap
}
