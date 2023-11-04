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
