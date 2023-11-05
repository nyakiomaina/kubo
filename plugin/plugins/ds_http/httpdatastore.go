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
	"github.com/ipfs/go-datastore/query"
)

type HttpDatastore struct {
	serverURL string
	client    *http.Client
}

func NewHttpDatastore(serverURL string) *HttpDatastore {
	fmt.Println("Initializing HttpDatastore for testing with httpbin...")
	return &HttpDatastore{
		serverURL: "https://httpbin.org",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *HttpDatastore) Batch(ctx context.Context) (ds.Batch, error) {
	return nil, errors.New("Batch not implemented")
}

func (h *HttpDatastore) Close() error {
	return errors.New("Close not implemented")
}

func (h *HttpDatastore) Delete(ctx context.Context, key ds.Key) error {
	return errors.New("Delete not implemented")
}

func (h *HttpDatastore) GetSize(ctx context.Context, key ds.Key) (size int, err error) {
	return -1, errors.New("GetSize not implemented")
}

func (h *HttpDatastore) Query(ctx context.Context, q query.Query) (query.Results, error) {
	return nil, errors.New("Query not implemented")
}

func (h *HttpDatastore) Sync(ctx context.Context, prefix ds.Key) error {
	return errors.New("Sync not implemented")
}

func (h *HttpDatastore) Put(ctx context.Context, key ds.Key, value []byte) error {
	fullURL := h.serverURL + "/anything/" + key.String()
	req, err := http.NewRequestWithContext(ctx, "PUT", fullURL, bytes.NewReader(value))
	if err != nil {
		return err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("server responded with an error: %v, body: %s", resp.Status, string(body))
	}

	fmt.Println("PUT request to:", fullURL)
	return nil
}

func (h *HttpDatastore) Get(ctx context.Context, key ds.Key) (value []byte, err error) {
	fullURL := h.serverURL + "/anything/" + key.String()
	resp, err := h.client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("server responded with an error: %v, body: %s", resp.Status, string(body))
	}

	fmt.Println("GET request to:", fullURL)
	return ioutil.ReadAll(resp.Body)
}

func (h *HttpDatastore) Has(ctx context.Context, key ds.Key) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "HEAD", h.serverURL+"/has/"+key.String(), nil)
	if err != nil {
		return false, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		return false, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}
}
