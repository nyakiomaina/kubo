package httpdatastore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/ipfs/boxo/datastore/dshelp"
	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

type HttpDatastore struct {
	serverURL string
	client    *http.Client
}

func NewHttpDatastore() *HttpDatastore {
	var serverURL string
	if envURL := os.Getenv("HTTP_DATASTORE_URL"); envURL != "" {
		serverURL = envURL
	} else {
		serverURL = "http://127.0.0.1:9500"
	}

	return &HttpDatastore{
		serverURL: serverURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *HttpDatastore) Batch(ctx context.Context) (ds.Batch, error) {
	return nil, errors.New("Batch not implemented")
}

func (h *HttpDatastore) Close() error {
	return nil
}

func (h *HttpDatastore) Delete(ctx context.Context, key ds.Key) error {
	return nil
}

func (h *HttpDatastore) GetSize(ctx context.Context, key ds.Key) (size int, err error) {
	return -1, nil
}

func (h *HttpDatastore) Query(ctx context.Context, q query.Query) (query.Results, error) {
	return nil, errors.New("Query not implemented")
}

func (h *HttpDatastore) Sync(ctx context.Context, prefix ds.Key) error {
	return nil
}

func (h *HttpDatastore) Push(ctx context.Context, prefix ds.Key) error {
	return nil
}

func (h *HttpDatastore) Put(ctx context.Context, key ds.Key, value []byte) error {
	fmt.Println("put")
	cidV1, err := dshelp.DsKeyToCidV1(key, cid.DagProtobuf)
	if err != nil {
		return fmt.Errorf("failed to convert key to CID: %v", err)
	}
	cidStr := cidV1.String()

	fullURL := h.serverURL + "/put/" + cidStr
	req, err := http.NewRequestWithContext(ctx, "PUT", fullURL, bytes.NewReader(value))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %v", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("server responded with an error: %v, body: %s", resp.Status, string(body))
	}

	fmt.Printf("PUT request successful for key: %s with CIDv1: %s\n", key, cidStr)
	return nil
}

func (h *HttpDatastore) Get(ctx context.Context, key ds.Key) (value []byte, err error) {
	fmt.Println("get")
	cidV1, err := dshelp.DsKeyToCidV1(key, cid.DagProtobuf)
	if err != nil {
		return nil, fmt.Errorf("failed to convert key to CID: %v", err)
	}
	cidStr := cidV1.String()

	fullURL := h.serverURL + "/get/" + cidStr
	resp, err := h.client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("server responded with an error: %v, body: %s", resp.Status, string(body))
	}

	fmt.Printf("GET request successful for key: %s with CIDv1: %s\n", key, cidStr)
	return ioutil.ReadAll(resp.Body)
}

// func (h *HttpDatastore) Has(ctx context.Context, key ds.Key) (bool, error) {
// 	return true, nil
// }

func (h *HttpDatastore) Has(ctx context.Context, key ds.Key) (bool, error) {
	cidV1, err := dshelp.DsKeyToCidV1(key, cid.DagProtobuf)
	if err != nil {
		return false, fmt.Errorf("failed to convert key to CID V1: %s", err)
	}

	fmt.Printf("Decoded CID V1: %s\n", cidV1)

	fullURL := h.serverURL + "/has/" + cidV1.String()
	req, err := http.NewRequestWithContext(ctx, "HEAD", fullURL, nil)
	if err != nil {
		return false, fmt.Errorf("creating HEAD request failed: %s", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("executing HEAD request failed: %s", err)
	}
	defer resp.Body.Close()
	exists := resp.StatusCode == http.StatusOK
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return false, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}
	fmt.Printf("CIDV1: %s exists: %t\n", cidV1, exists)

	return exists, nil
}
