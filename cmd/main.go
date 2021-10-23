package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nxadm/tail"
	"github.com/rahulkhairwar/logtail"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	file string
)

var config logtail.Config

func setup() {
	file = os.Getenv("FILE_TO_TAIL")

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

	conf := tail.Config{
		MustExist: true,
		Poll:      true,
		Follow:    true,
	}
	t, err := tail.TailFile(file, conf)
	if err != nil {
		log.Fatalf("failed to tail file {%v}, err: %v", file, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}

	wg.Add(1)
	defer wg.Wait()

	// start HTTP server in a separate goroutine
	go logtail.Serve(ctx, &config, &wg)

	for {
		select {
		case line := <-t.Lines:
			logtail.HandleLogLine(&config, line.Text)
		case <-t.Dead():
			break
		}
	}
}
