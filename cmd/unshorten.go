package main

import (
	"context"
	"flag"
	"github.com/sfomuseum/go-url-unshortener"
	"log"
	"time"
)

func main() {

	stdin := flag.Bool("stdin", false, "Read URLs from STDIN")

	flag.Parse()

	rate := time.Second / 10

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

	remaining := 0

	done_ch := make(chan bool)
	err_ch := make(chan error)

	unshorten := func(ctx context.Context, str_url string) {

		defer func() {
			done_ch <- true
		}()

		u, err := unshortener.UnshortenString(ctx, cache, str_url)

		if err != nil {
			err_ch <- err
			return
		}

		log.Println(u.String())
	}

	if *stdin {

	} else {

		for _, str_url := range flag.Args() {
			remaining += 1
			go unshorten(ctx, str_url)
		}
	}

	for remaining > 0 {

		select {
		case <-done_ch:
			remaining -= 1
		case err := <-err_ch:
			log.Println(err)
		default:
			// pass
		}
	}
}
