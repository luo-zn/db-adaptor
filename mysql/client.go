
package mysql

import (
	"context"
	"github.com/luo-zn/db-adaptor/bases"
)

type MysqlClient struct {
	client string
	Ctx    context.Context
}

func NewMysqlClient(uri string) *MysqlClient {
	return &MysqlClient{client: uri}
}

func (MysqlClient) Create(tb string, entity bases.Entity) (interface{}, error) {
	panic("implement me")
}

func (MysqlClient) Retrieve(tb string, filter bases.Entity) (e interface{}, err error) {
	panic("implement me")
}

func (MysqlClient) Update(tb string, e bases.Entity) (bool, error) {
	panic("implement me")
}

func (MysqlClient) Delete(tb string, e bases.Entity) (bool, error) {
	panic("implement me")
}
