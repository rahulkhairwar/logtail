package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/nxadm/tail"
	"github.com/rahulkhairwar/logtail"
	"log"
	"os"
)

var (
	file string
)

func init() {
	file = os.Getenv("TAIL_FILE")
}

func main() {
	conf := tail.Config{
		MustExist: true,
		Poll:      true,
		Follow:    true,
	}
	t, err := tail.TailFile(file, conf)
	if err != nil {
		log.Fatalf("failed to tail file {%v}, err: %v", file, err)
	}

	for {
		select {
		case line := <-t.Lines:
			// log.Println(line.Text)
			logtail.HandleLogLine(line.Text)
		case <-t.Dead():
			break
		}
	}
}
