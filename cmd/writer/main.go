package main

import (
	"context"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nxadm/tail"
	"log"
	"os"
	"sync"
	"time"
)

var (
	inFile, outFile string
)

func init() {
	inFile = os.Getenv("SAMPLE_LOG_FILE")
	outFile = os.Getenv("OUTPUT_LOG_FILE")
}

const (
	stopCount = 1e3
	lifetime  = 1 * time.Minute
)

func main() {
	ctx := context.Background()
	conf := tail.Config{
		MustExist: true,
		Follow:    false,
	}
	in, err := tail.TailFile(inFile, conf)
	if err != nil {
		log.Fatalf("failed to tail reader outFile {%v}, err: %v", inFile, err)
	}
	defer in.Cleanup()

	out, err := os.OpenFile(outFile, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalf("failed to open outFile {%v}, err: %v", outFile, err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Fatalf("failed to close out outFile {%v}, err: %v", outFile, err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()

	lineCnt := 0
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go regulateLifetime(ctx, in, lifetime, &wg)

loop:
	for {
		select {
		case line := <-in.Lines:
			if line == nil {
				continue
			}

			line.Text += "\n"

			if _, err := out.Write([]byte(line.Text)); err != nil {
				log.Printf("[ERROR] failed to write line text to file, err: %v\n", err)
			}

			lineCnt++
			if lineCnt == stopCount {
				fmt.Printf("killing after reading {%v} lines\n", stopCount)
				in.Kill(fmt.Errorf("read {%v} lines", stopCount))
			}

			// sleep for 1 second after writing each line to emulate live logs.
			time.Sleep(time.Second)
		case <-in.Dead():
			break loop
		}
	}
}

func regulateLifetime(ctx context.Context, t *tail.Tail, d time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	if ctx == nil || t == nil {
		return
	}

	select {
	case <-ctx.Done():
		fmt.Printf("ctx done, ctx.Err: %v\n", ctx.Err())
		return
	case <-time.After(d):
		fmt.Printf("[regulateLifetime] maximum allowed lifetime of {%v} completed\n", d)
		t.Kill(nil)
		return
	}
}
