package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-url-unshortener"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	progress := flag.Bool("progress", false, "Display progress information")
	verbose := flag.Bool("verbose", false, "Be chatty about what's going on")
	stdin := flag.Bool("stdin", false, "Read URLs from STDIN")
	qps := flag.Int("qps", 10, "Number of (unshortening) queries per second")

	flag.Parse()

	rate := time.Second / time.Duration(*qps)

	worker, err := unshortener.NewThrottledUnshortener(rate)

	if err != nil {
		log.Fatal(err)
	}

	cache, err := unshortener.NewCachedUnshortener(worker)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signal_ch := make(chan os.Signal)
	signal.Notify(signal_ch, os.Interrupt, syscall.SIGTERM)

	go func(c chan os.Signal) {
		
		<-c

		cancel()
		os.Exit(0)
	}(signal_ch)

	remaining := 0

	type UnshortenResponse struct {
		ShortenedURL   string
		UnshortenedURL string
		Error          error
	}

	done_ch := make(chan bool)
	rsp_ch := make(chan *UnshortenResponse)

	unshorten := func(ctx context.Context, str_url string) {

		defer func() {
			done_ch <- true
		}()

		u, err := unshortener.UnshortenString(ctx, cache, str_url)

		var rsp UnshortenResponse

		if err != nil {

			rsp = UnshortenResponse{
				ShortenedURL: str_url,
				Error:        err,
			}

			rsp_ch <- &rsp
		}
		
		if u != nil {

			rsp = UnshortenResponse{
				ShortenedURL:   str_url,
				UnshortenedURL: u.String(),
			}

			rsp_ch <- &rsp
		}

		// assume that ctx.Done() has been invoked
	}

	if *stdin {

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			remaining += 1
			str_url := scanner.Text()
			go unshorten(ctx, str_url)
		}

	} else {

		for _, str_url := range flag.Args() {
			remaining += 1
			go unshorten(ctx, str_url)
		}
	}

	total := remaining

	completed_ch := make(chan bool)

	if *progress {

		go func() {

			for {
				select {
				case <-completed_ch:
					break
				case <-time.After(1 * time.Minute):
					log.Printf("%d of %d URL left to unshorten\n", remaining, total)
				}
			}
		}()
	}

	lookup := new(sync.Map)

	for remaining > 0 {

		select {
		case <-done_ch:
			remaining -= 1
		case rsp := <-rsp_ch:

			if rsp.Error != nil {

				lookup.Store(rsp.ShortenedURL, "-")
				log.Printf("Failed to unshorted '%s' %s", rsp.ShortenedURL, rsp.Error)

			} else {

				lookup.Store(rsp.ShortenedURL, rsp.UnshortenedURL)

				if *verbose {
					log.Printf("%s becomes %s\n", rsp.ShortenedURL, rsp.UnshortenedURL)
				}
			}

		default:
			// log.Println(remaining)
		}
	}

	completed_ch <- true

	report := make(map[string]string)

	lookup.Range(func(k interface{}, v interface{}) bool {
		shortened_url := k.(string)
		unshortened_url := v.(string)
		report[shortened_url] = unshortened_url
		return true
	})

	enc_report, err := json.Marshal(report)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(enc_report))
}
