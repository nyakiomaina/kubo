package httpdatastore

import (
	"fmt"

	"github.com/ipfs/kubo/plugin"
	"github.com/ipfs/kubo/repo/fsrepo"
)

type HttpPlugin struct{}

func (p *HttpPlugin) Name() string {
	return "ds_http"
}

func (p *HttpPlugin) Version() string {
	return "0.1.0"
}

func (p *HttpPlugin) Init(env *plugin.Environment) error {
	fmt.Println("plugin loaded!!")
	return nil
}

func (p *HttpPlugin) DatastoreTypeName() string {
	return "ds_http"
}

func (p *HttpPlugin) DatastoreConfigParser() fsrepo.ConfigFromMap {
	return func(m map[string]interface{}) (fsrepo.DatastoreConfig, error) {
		cfg := &HttpConfig{}
		err := cfg.ConfigFromMap(m)
		return cfg, err
	}
}

var Plugins = []plugin.Plugin{
	&HttpPlugin{},
}

// Comment out or remove these redefined methods and struct from this file

// type HttpConfig struct {
// 	// ... your fields
// }

// var _ fsrepo.DatastoreConfig = (*HttpConfig)(nil)

// func (cfg *HttpConfig) DiskSpec() fsrepo.DiskSpec {
// 	...
// }

// func (cfg *HttpConfig) Create(path string) (repo.Datastore, error) {
// 	...
// }

// func (cfg *HttpConfig) ConfigFromMap(m map[string]interface{}) error {
// 	...
// }
