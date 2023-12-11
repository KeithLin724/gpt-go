package pkg

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type FetchResult struct {
	URL     string
	State   string
	Connect bool

	mx sync.Mutex
}

func (ptr *FetchResult) SetData(url, state string, connect bool) {
	ptr.mx.Lock()

	ptr.URL = url
	ptr.State = state
	ptr.Connect = connect

	ptr.mx.Unlock()
}
func (ptr *FetchResult) SetURL(url string) {
	ptr.mx.Lock()

	ptr.URL = url

	ptr.mx.Unlock()
}

func (ptr *FetchResult) GetURL() (res string) {
	ptr.mx.Lock()

	res = ptr.URL

	ptr.mx.Unlock()
	return
}

func (ptr *FetchResult) GetState() (res string) {
	ptr.mx.Lock()

	res = ptr.State

	ptr.mx.Unlock()
	return
}

func (ptr *FetchResult) GetConnect() (res bool) {
	ptr.mx.Lock()

	res = ptr.Connect

	ptr.mx.Unlock()
	return
}

// The `FetchHTTP` function is a method of the `FetchResult` struct. It is responsible for fetching the
// HTTP response from the specified URL and updating the state of the `FetchResult` instance.
func (ptr *FetchResult) FetchHTTP(log *Log) {
	url, state, connect := ptr.URL, "Success", true

	resp, err := http.Get(url)

	if err != nil {
		state = fmt.Sprintf("Error: %s", err)
		connect = false
		log.Error(state)
	}

	ptr.SetData(url, state, connect)

	if connect {
		log.Info("Success")
		defer resp.Body.Close()
	}
}

// The `RunFetchServer` function is a method of the `FetchResult` struct. It starts a new goroutine
// that calls the `FetchHTTP` method of the `FetchResult` instance.
func (ptr *FetchResult) RunFetchServer(sec int64) {
	log := NewLog()
	go func() {
		for {
			ptr.FetchHTTP(log)
			time.Sleep(time.Duration(sec) * time.Second)
		}
	}()
}
