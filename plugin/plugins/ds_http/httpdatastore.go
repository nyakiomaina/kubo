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
	dsq "github.com/ipfs/go-datastore/query"
)

type HttpDatastore struct {
	serverURL string
	client    *http.Client

	localDatastore map[ds.Key][]byte
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
		localDatastore: make(map[ds.Key][]byte),
	}
}

func (s *HttpDatastore) Batch(_ context.Context) (ds.Batch, error) {
	return nil, nil
}
func (s *HttpDatastore) Close() error {
	return nil
}

func (s *HttpDatastore) GetSize(ctx context.Context, k ds.Key) (size int, err error) {
	return 0, nil
}

func (s *HttpDatastore) Query(ctx context.Context, q dsq.Query) (dsq.Results, error) {
	return nil, nil
}
// Put sends the given value to the HTTP server to be stored under the specified key.
func (h *HttpDatastore) Put(ctx context.Context, key ds.Key, value []byte) error {
	req, err := http.NewRequestWithContext(ctx, "PUT", h.serverURL+"/put/"+key.String(), bytes.NewReader(value))
	if err != nil {
		return err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("server responded with an error")
	}
	return nil
}
func (s *HttpDatastore) Sync(ctx context.Context, prefix ds.Key) error {
	return nil
}

func (h *HttpDatastore) Get(ctx context.Context, key ds.Key) (value []byte, err error) {
	resp, err := h.client.Get(h.serverURL + "/get/" + key.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("server responded with an error")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
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

// fetch data from IPFS given the CID
func (h *HttpDatastore) GetByCID(ctx context.Context, cid string) ([]byte, error) {
	resp, err := h.client.Get(h.serverURL + "/ipfs/" + cid)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to GET data from IPFS")
	}

	return ioutil.ReadAll(resp.Body)
}

// upload data to IPFS and return the CID of the uploaded file
func (h *HttpDatastore) PostState(ctx context.Context, state []byte) (string, error) {
	req, err := http.NewRequest("POST", h.serverURL+"/ipfs/", bytes.NewReader(state))
	if err != nil {
		return "", err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to POST state to IPFS")
	}

	//read the CID from the response
	cid, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(cid), nil
}
