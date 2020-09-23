package main

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/jeremywohl/flatten"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
)

func NewIgnoreItems(gdIgnorePath string) map[string][]string {
	file, err := ioutil.ReadFile(gdIgnorePath)
	if err != nil {
		log.Fatalf("Unable to open gdignore.yaml: %v", err)
	}

	jsonContent, err := yaml.YAMLToJSON(file)
	if err != nil {
		log.Errorf("Unable to convert yaml to json: %v", err)
	}

	flatConfigString, err := flatten.FlattenString(string(jsonContent), "", flatten.PathStyle)
	if err != nil {
		log.Errorf("Unable to parse gdignore.yaml to JSON string: %v", err)
	}

	var flatConfig map[string]string
	if err := json.Unmarshal([]byte(flatConfigString), &flatConfig); err != nil {
		log.Errorf("Unable to parse json: %v", err)
	}

	return toListValues(flatConfig)
}

func toListValues(config map[string]string) map[string][]string {
	var pathItems = make(map[string][]string)

	for key := range config {
		path := keyToPath(key)
		appendToMap(pathItems, path, config[key])
	}

	return pathItems
}

func keyToPath(key string) string {
	items := strings.Split(key, PathSep)

	var tokens []string
	for i := 0; i < len(items); i = i + 2 {
		tokens = append(tokens, items[i])
	}
	return strings.Join(tokens, PathSep)
}
