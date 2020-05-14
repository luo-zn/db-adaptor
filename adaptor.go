package db_adaptor

import (
	"github.com/luo-zn/db-adaptor/mongodb"
	"github.com/luo-zn/db-adaptor/mysql"
)

type DbAdaptor struct {
	DbC DBClient
	Opt *AdaptorOptions
}

func NewDbAdaptor(opts ...*AdaptorOptions) *DbAdaptor {
	var adaptor DbAdaptor
	clientOpt := MergeAdaptorOptions(opts...)
	adaptor.Opt = clientOpt
	adaptor.newClient()
	return &adaptor
}

func (db *DbAdaptor) newClient() {
	if db.Opt.DBType == "mysql" {
		db.DbC = mysql.NewMysqlClient(db.Opt.Uri)
	} else {
		db.Opt.DBType = "mongodb"
		mgC := mongodb.NewMgoClient(db.Opt.Uri)
		db.DbC = mgC
	}
}

//func NewClient(uri string) *Client {
//	c, _ := mongodb.NewMongoClient(uri)
//	return &Client{ DbClient: c}
//}
