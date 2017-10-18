package actions

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func actionsPath(b *backend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "actions/?",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: b.actionsPathList,
			},
		},
		&framework.Path{
			Pattern: "actions/" + framework.GenericNameRegex("action"),
			Fields: map[string]*framework.FieldSchema{
				"key": &framework.FieldSchema{Type: framework.TypeString},
			},
			ExistenceCheck: b.pathExistenceCheck,
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: b.actionsCreateUpdate,
			},
		},
	}
}

func (b *backend) pathExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %v", err)
	}

	return out != nil, nil
}

func (b *backend) actionsPathList(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	vals, err := req.Storage.List("actions/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(vals), nil
}

func (b *backend) actionsCreateUpdate(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry := &logical.StorageEntry{
		Key:   req.Path,
		Value: []byte("1"),
	}
	s := req.Storage
	err := s.Put(entry)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"value": 1,
		},
	}, nil
}
