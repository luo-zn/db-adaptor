package mysql

import (
	"context"
	"github.com/luo-zn/db-adaptor/bases"
)

type MysqlClient struct {
	client string
	Ctx    context.Context
}

func (MysqlClient) Count(tb string, filter bases.Entity) (int64, error) {
	panic("implement me")
}

func (MysqlClient) Connect(opt map[string]interface{}) error {
	panic("implement me")
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

func (MysqlClient) Update(tb string, f bases.Entity, e bases.Entity) (bool, error) {
	panic("implement me")
}

func (MysqlClient) UpdateOneWithFilter(tb string, filter map[string]interface{}, e bases.Entity) (bool, error) {
	panic("implement me")
}

func (MysqlClient) Delete(tb string, e bases.Entity) (bool, error) {
	panic("implement me")
}

func (MysqlClient) Close() error {
	panic("implement me")
}
