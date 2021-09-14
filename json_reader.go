package logtail

import (
	"encoding/json"
	"log"
)

const (
	separator = "-----------------------------------------------------------------------------------------"
)

func HandleLogLine(conf *Config, l string) {
	m, err := stringToJSON(l)
	if err != nil {
		log.Println("failed to convert string to JSON, err: ", err)

		return
	}

	printSelectedKeys(conf, m)
}

func printSelectedKeys(conf *Config, m map[string]interface{}) {
	for k, v := range m {
		if contains(conf.ParseKeys, k) {
			log.Printf("%v = %v\n", k, v)
		}
	}

	log.Println(separator)
}

func contains(arr []string, key string) bool {
	for _, s := range arr {
		if key == s {
			return true
		}
	}

	return false
}

func stringToJSON(s string) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, err
	}

	return m, nil
}
