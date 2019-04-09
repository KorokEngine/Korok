// +build js

package res

import (
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
