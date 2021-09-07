package logtail

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

const (
	separator = "-----------------------------------------------------------------------------------------"
)

var conf config

func init() {
	f, err := os.OpenFile("./config.yaml", os.O_RDONLY, 0755)
	if err != nil {
		log.Fatalf("failed to open config file, err: %v", err)
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("failed to read config file, err: %v", err)
	}

	// todo: can possibly implement UnmarshalYAML to get behaviour of reading
	//  ["File", "Function", "Level", "Line", "Message", "Time"] into a map[string]bool.
	if err := yaml.Unmarshal(bytes, &conf); err != nil {
		log.Fatalf("failed to unmarshal config, err: %v", err)
	}

	log.Printf("parsed config: %v\n", conf)
}

func HandleLogLine(l string) {
	m, err := stringToJSON(l)
	if err != nil {
		log.Println("failed to convert string to JSON, err: ", err)
		return
	}

	printSelectedKeys(&conf, m)
}

func printSelectedKeys(conf *config, m map[string]interface{}) {
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
