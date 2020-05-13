/* Create By LZN */
package mongodb

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

var (
	dbUri = "mongodb://localhost:27017"
	uid   = "5eb8e1f6f715a2494825d0ba"
)

type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Nickname string `json:"nickname,omitempty" bson:"nickname,omitempty"`
}

func (u *User) DataBase() string {
	return "testDB"
}

func TestMgoClient_Connect(t *testing.T) {
	mg := &MgoClient{}
	mg.Connect(dbUri, 20)
	err := mg.client.Ping(context.TODO(), readpref.Primary())
	assert.Nil(t, err)
	assert.NotEmpty(t, mg.client)
}

func TestMgoClient_getCollection(t *testing.T) {
	mg := &MgoClient{}
	mg.Connect(dbUri, 20)
	col := mg.getCollection("testDB", "user")
	assert.NotNil(t, col)
	assert.Equal(t, col.Name(), "user")
}

func TestMgoClient_FindOne(t *testing.T) {
	u := &User{Username: "testUser", Nickname: "testNickname"}
	mg := &MgoClient{}
	mg.Connect(dbUri, 20)
	_, err := mg.FindOne("user", u)
	assert.Nil(t, err, "Error %s", err)
	assert.Equal(t, "testUser", u.Username, "The expected Username is testUser")
	assert.Equal(t, "testNickname", u.Nickname, "The expected Nickname is testNickname")
	assert.NotEmpty(t, u.ID, "User id is %s", u.ID)
}

func TestMgoClient_Create(t *testing.T) {
	u := &User{Username: "testUser", Nickname: "testNickname", ID: uid}
	mg := &MgoClient{}
	mg.Connect(dbUri, 20)
	res, err := mg.Create("user", u)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMgoClient_Retrieve(t *testing.T) {
	u := &User{Username: "testUser", Nickname: "testNickname"}
	mg := &MgoClient{}
	mg.Connect(dbUri, 20)
	mg.FindOne("user", u)
	assert.Equal(t, u.Username, "testUser", "The expected Username is testUser")
	assert.Equal(t, u.Nickname, "testNickname", "The expected Nickname is testNickname")
	assert.NotEmpty(t, u.ID, "User id is %s", u.ID)
	assert.Equal(t, uid, u.ID, "The expected ID is %s", uid)
}

func TestMgoClient_UpdateOneById(t *testing.T) {
	u := &User{Username: "testUser", Nickname: "testNickname"}
	mg := &MgoClient{}
	mg.Connect(dbUri, 20)
	res, err := mg.UpdateOneById("user", uid, u)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestMgoClient_Update(t *testing.T) {
	u := &User{Username: "testUser", Nickname: "testNickname", ID: uid}
	mg := &MgoClient{}
	mg.Connect(dbUri, 20)
	res, err := mg.Update("user", u)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestMgoClient_Delete(t *testing.T) {
	u := &User{Username: "testUser", Nickname: "testNickname", ID: uid}
	mg := &MgoClient{}
	mg.Connect(dbUri, 20)
	res, err := mg.Delete("user", u)
	assert.Nil(t, err)
	assert.True(t, res)
}
