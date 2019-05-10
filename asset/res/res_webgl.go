// +build js

package res

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

// Open opens a named asset.
func Open(name string) (File, error) {

	client := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(name)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("resp.StatusCode=%d", resp.StatusCode)
	}
	return resp.Body, nil
}

type wfile struct {
	name string
}

func (wf *wfile) Write(p []byte) (n int, err error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Post(wf.name, "", bytes.NewReader(p))
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return 0, fmt.Errorf("resp.StatusCode=%d", resp.StatusCode)
	}
	return len(p), nil
}

// Open opens a named asset.
func Create(name string) (WFile, error) {
	return &wfile{name: name}, nil
}
