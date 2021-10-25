package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rahulkhairwar/logtail"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

var config logtail.Config

func setup() {
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
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		log.Fatalf("failed to unmarshal config, err: %v", err)
	}

	log.Printf("parsed config: %v\n", config)
}

func main() {
	setup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start the HTTP server.
	logtail.Serve(ctx, &config)
}
