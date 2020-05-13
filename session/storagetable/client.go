package client

import (
	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/jeffguorg/middlewares/session"
)

type Client struct {
	table     *storage.Table
	getEntity func(table *storage.Table, sessionId string) *storage.Entity
}

// NewConnectionStringClient returns a client for session usage
func NewConnectionStringClient(connstr, name string) (*Client, error) {
	cli, err := storage.NewClientFromConnectionString(connstr)
	if err != nil {
		return nil, err
	}
	tblSrv := cli.GetTableService()
	return &Client{
		table: tblSrv.GetTableReference(name),
	}, nil
}

func (client Client) Load(name string) (map[string]interface{}, error) {
	var entity *storage.Entity
	if client.getEntity != nil {
		entity = client.getEntity(client.table, name)
	} else {
		entity = client.table.GetEntityReference("session", name)
	}
	if err := entity.Get(10, storage.MinimalMetadata, nil); err != nil {
		return nil, err
	}
	return entity.Properties, nil
}

func (client Client) Reset(name string) error {
	var entity *storage.Entity
	if client.getEntity != nil {
		entity = client.getEntity(client.table, name)
	} else {
		entity = client.table.GetEntityReference("session", name)
	}
	if err := entity.Get(10, storage.MinimalMetadata, nil); err != nil {
		if err := entity.Delete(false, nil); err != nil {
			return err
		}
	}
	return nil
}

func (client Client) Update(name string, value map[string]interface{}) error {
	var entity *storage.Entity
	if client.getEntity != nil {
		entity = client.getEntity(client.table, name)
	} else {
		entity = client.table.GetEntityReference("session", name)
	}
	entity.Properties = value
	if err := entity.InsertOrMerge(nil); err != nil {
		return err
	}
	return nil
}

var (
	_ session.Client = Client{}
)
