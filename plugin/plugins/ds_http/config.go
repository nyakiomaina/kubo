package httpdatastore

import (
	"github.com/ipfs/kubo/repo"
	"github.com/ipfs/kubo/repo/fsrepo"
)

type HttpConfig struct {
	ServerURL string
}

func (cfg *HttpConfig) DiskSpec() fsrepo.DiskSpec {
	return fsrepo.DiskSpec{
		"type":      "httpDatastore",
		"serverURL": cfg.ServerURL,
	}
}

func (cfg *HttpConfig) ConfigFromMap(m map[string]interface{}) error {
	// "serverURL"??
	if url, ok := m["serverURL"]; ok {
		cfg.ServerURL = url.(string)
	}
	return nil
}

func (cfg *HttpConfig) Create(path string) (repo.Datastore, error) {
	return NewHttpDatastore(cfg.ServerURL), nil
}
