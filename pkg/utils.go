package pkg

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func GetConfig[Config any](configPath *string) Config {
	var conf Config
	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	if err = yaml.Unmarshal(data, &conf); err != nil {
		log.Fatal(err)
	}
	return conf
}
