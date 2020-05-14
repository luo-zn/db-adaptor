# db-adaptor
db-adaptor使用例子

## 定义User model
```go
package models

import "github.com/luo-zn/db-adaptor"

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