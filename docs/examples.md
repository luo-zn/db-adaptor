# db-adaptor
db-adaptor使用例子

## 定义global
```go
package global

import (
	"fmt"
	"github.com/luo-zn/db-adaptor"
)

var (
	//DBUri Database connect string.
	DBUri string
	//DBAdapt Database Adaptor
	DBAdapt *db_adaptor.DbAdaptor
)

func init() {
	fmt.Print("Call global init.")
}

func Init(dbUri string)  {
	DBUri = dbUri
	opt := db_adaptor.AdaptorOptions{Uri: DBUri,DBType:"mongodb"}
	DBAdapt = db_adaptor.NewDbAdaptor(&opt)
	// or
	//opt := db_adaptor.AdaptorOptions{Uri: DBUri}
    //opt1 := db_adaptor.AdaptorOptions{DBType:"mongodb"}
    //DBAdapt = db_adaptor.NewDbAdaptor(&opt,&opt1)
}

```
## 定义User model
```go
package models

import "models/global"

// baseModel  defined common properties.
type baseModel struct {
	tableName string
	dataBase  string
}
// User need to implement DBClient interface.
type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Nickname string `json:"nickname,omitempty" bson:"nickname,omitempty"`
	
	baseModel `bson:"-"`
}
// DataBase getDataBase name.
// implement bases.Entity interface
func (u *User) DataBase() string {
	return u.dataBase
}
//Create insert a new data to Database.
func (u *User) Create() error {
	res, err := global.DBAdapt.DbC.Create(u.tableName, u)
	u.ID = res.(string)
	return err
}
//Get find data from DataBase.
func (u *User) Get() error {
	panic("implement me")
}
//Update instance property values to Database.
func (u *User) Update() (bool, error) {
	panic("implement me")
}
// Delete will remove a data from Database.
func (u *User) Delete() (bool, error) {
	panic("implement me")
}
```