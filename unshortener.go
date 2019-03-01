package unshortener

import (
	_ "log"
	"net/http"
	"net/url"
	"time"
	"sync"
)

type Unshortener interface {
	UnshortenString(u string) (*url.URL, error)
	Unshorten(u *url.URL) (*url.URL, error)
}

type ThrottledUnshortener struct {
	Unshortener
	throttle <-chan time.Time
	client   *http.Client
}

type CachedUnshortener struct {
	Unshortener
	worker Unshortener
	cache *sync.Map
}

func UnshortenString(un Unshortener, str_u string) (*url.URL, error) {

	u, err := url.Parse(str_u)

	if err != nil {
		return nil, err
	}

	return un.Unshorten(u)
}

func NewCachedUnshortener(worker Unshortener) (Unshortener, error) {

	cache := new(sync.Map)
	
	un := CachedUnshortener{
		worker: worker,
		cache: cache,
	}

	return &un, nil
}

func (un *CachedUnshortener) Unshorten(u *url.URL) (*url.URL, error) {

	str_url := u.String()
	
	v, ok := un.cache.Load(str_url)

	if ok {
		str_url = v.(string)
		return url.Parse(str_url)
	}

	u2, err := un.worker.Unshorten(u)

	if err != nil {
		return nil, err
	}

	un.cache.Store(u.String(), u2.String())
	return u2, nil
}

func NewThrottledUnshortener(rate time.Duration) (Unshortener, error) {

	throttle := time.Tick(rate)

	client := &http.Client{
		// something something something client.CheckRedirect - configure for more than (default number of) hops?
		// https://stackoverflow.com/questions/23297520/how-can-i-make-the-go-http-client-not-follow-redirects-automatically
		// https://jonathanmh.com/tracing-preventing-http-redirects-golang/
	}

	un := ThrottledUnshortener{
		throttle: throttle,
		client:   client,
	}

	return &un, nil
}

func (un *ThrottledUnshortener) Unshorten(u *url.URL) (*url.URL, error) {

	<-un.throttle

	rsp, err := un.client.Head(u.String())

	if err != nil {
		return nil, err
	}

	return rsp.Request.URL, nil
}