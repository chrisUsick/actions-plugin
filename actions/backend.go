package actions

import (
	"log"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/plugins"
)

//New returns a new backend as an interface
func New() (interface{}, error) {
	return Backend(), nil
}
func Run(apiTLSConfig *api.TLSConfig) {
	actions, err := New()
	if err != nil {
		panic(err)
	}
	plugins.Serve(actions.(backend), apiTLSConfig)
}

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	log.Println("In the factory")
	b := Backend()
	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help:  "",
		Paths: framework.PathAppend(actionsPath(&b)),

		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}

	return &b
}

type backend struct {
	*framework.Backend
}
