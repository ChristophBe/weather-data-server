package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"sync"
)

type ConfigurationManager struct {
	configuration *Configuration
}

func (c *ConfigurationManager) LoadConfig(configFile string) (err error) {

	log.Printf("Loading Configuration Configfile %s", configFile)

	jsonFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Printf("failed to read configuration error: %v ", err)
		return
	}

	err = json.Unmarshal(jsonFile, &c.configuration)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return
}
func (c *ConfigurationManager) GetConfig() (config *Configuration, err error) {

	if c.configuration == nil {
		err = errors.New("configuration is not loaded")
	}
	config = c.configuration
	return
}

var mut sync.Mutex

var confManager *ConfigurationManager

func GetConfigManager() *ConfigurationManager {
	mut.Lock()
	defer mut.Unlock()

	if confManager == nil {
		confManager = &ConfigurationManager{}
	}
	return confManager
}
