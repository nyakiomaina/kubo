package httpdatastore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	ds "github.com/ipfs/go-datastore"
)

type HttpDatastore struct {
	serverURL string
	client    *http.Client
}

// func NewHttpDatastore(serverURL string) *HttpDatastore {
// 	fmt.Println("Initializing HttpDatastore for IPFS daemon...")
// 	return &HttpDatastore{
// 		serverURL: serverURL,
// 		client: &http.Client{
// 			Timeout: 30 * time.Second,
// 		},
// 	}
// }

func NewHttpDatastore() *HttpDatastore {
	fmt.Println("Initializing HttpDatastore for testing with httpbin...")
	return &HttpDatastore{
		serverURL: "https://httpbin.org",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (h *HttpDatastore) Put(ctx context.Context, key ds.Key, value []byte) error {
	// API provided by cartesi machine? or communicate with it
	req, err := http.NewRequest("PUT", h.serverURL+"/put/"+key.String(), bytes.NewReader(value))
	if err != nil {
		return err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Failed to PUT data to the server")
	}

	return nil
}

func (h *HttpDatastore) Get(ctx context.Context, key ds.Key) (value []byte, err error) {

	resp, err := h.client.Get(h.serverURL + "/get/" + key.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ds.ErrNotFound
	}

	return ioutil.ReadAll(resp.Body)
}

func (h *HttpDatastore) Delete(ctx context.Context, key ds.Key) error {

	req, err := http.NewRequest("DELETE", h.serverURL+"/delete/"+key.String(), nil)
	if err != nil {
		return err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Failed to DELETE data from the server")
	}

	return nil
}

func (h *HttpDatastore) Has(ctx context.Context, key ds.Key) (exists bool, err error) {
	return true, nil
}
