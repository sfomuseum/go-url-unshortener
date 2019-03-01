package main

import (
	"flag"
	"log"
	"github.com/sfomuseum/go-url-unshortener"
	"time"
)

func main() {

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
	
	for _, str_url := range flag.Args(){

		u, err := unshortener.UnshortenString(cache, str_url)

		if err != nil {
			log.Fatal(err)
		}

		log.Println(u)
	}
	
}
