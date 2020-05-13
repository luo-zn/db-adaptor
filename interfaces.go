/* Create By Jenner.luo */
package db_adaptor

import "github.com/luo-zn/db-adaptor/bases"

//DBClient
type DBClient interface {
	Create(tb string, entity bases.Entity) (interface{}, error)
	Retrieve(tb string, filter bases.Entity) (interface{}, error)
	Update(tb string, e bases.Entity) (bool, error)
	Delete(tb string, e bases.Entity) (bool, error)
}
