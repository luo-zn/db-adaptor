/* Create By LZN */
package mongodb

import (
	"bou.ke/monkey"
	"context"
	"github.com/luo-zn/db-adaptor/bases"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"testing"
	"time"
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

func getMgoClient() *MgoClient {
	opt := map[string]interface{}{"uri": dbUri, "ctx_timeout": 40 * time.Second}
	mg := &MgoClient{}
	mg.Connect(opt)
	return mg
}

func TestMgoClient_Connect(t *testing.T) {
	var mgc *mongo.Client
	guard := monkey.PatchInstanceMethod(reflect.TypeOf(mgc),"Ping",
		func(_ *mongo.Client, ctx context.Context, rp *readpref.ReadPref) error{
			return nil
	})
	defer guard.Unpatch()
	mg := getMgoClient()
	err := mg.client.Ping(context.TODO(), readpref.SecondaryPreferred())
	assert.Nil(t, err)
	assert.NotEmpty(t, mg.client)
	assert.Nil(t, mg.Close())
}

func TestMgoClient_getCollection(t *testing.T) {
	mg := getMgoClient()
	col := mg.getCollection("testDB", "user")
	assert.NotNil(t, col)
	assert.Equal(t, col.Name(), "user")
}

func TestMgoClient_Create(t *testing.T) {
	var mgc *mongo.Collection
	guard := monkey.PatchInstanceMethod(reflect.TypeOf(mgc),"InsertOne",
		func(_ *mongo.Collection, ctx context.Context, document interface{},opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error){
			return &mongo.InsertOneResult{InsertedID:uid}, nil
		})
	defer guard.Unpatch()
	u := &User{Username: "testUser", Nickname: "testNickname", ID: uid}
	mg := getMgoClient()
	res, err := mg.Create("user", u)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestMgoClient_FindOne(t *testing.T) {
	var mgc *mongo.SingleResult
	guard := monkey.PatchInstanceMethod(reflect.TypeOf(mgc),"Decode",
		func(_ *mongo.SingleResult, v interface{}) error{
			obj, ok :=v.(User)
			if ok{
				obj.ID = uid
				return nil
			}
			return bases.Error("Obj is not User struct!")
		})
	defer guard.Unpatch()
	u := &User{Username: "testUser", Nickname: "testNickname"}
	mg := getMgoClient()
	_, err := mg.FindOne("user", u)
	assert.Nil(t, err, "Error %s", err)
	assert.Equal(t, "testUser", u.Username, "The expected Username is testUser")
	assert.Equal(t, "testNickname", u.Nickname, "The expected Nickname is testNickname")
	assert.NotEmpty(t, u.ID, "User id is %s", u.ID)
}

func TestMgoClient_Retrieve(t *testing.T) {
	u := &User{Username: "testUser", Nickname: "testNickname"}
	mg := getMgoClient()
	mg.FindOne("user", u)
	assert.Equal(t, u.Username, "testUser", "The expected Username is testUser")
	assert.Equal(t, u.Nickname, "testNickname", "The expected Nickname is testNickname")
	assert.NotEmpty(t, u.ID, "User id is %s", u.ID)
	assert.Equal(t, uid, u.ID, "The expected ID is %s", uid)
}

func TestMgoClient_UpdateOneById(t *testing.T) {
	u := &User{Username: "testUser", Nickname: "testNickname"}
	mg := getMgoClient()
	res, err := mg.UpdateOneById("user", uid, u)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestMgoClient_Update(t *testing.T) {
	u := &User{Username: "testUser", Nickname: "testNickname", ID: uid}
	mg := getMgoClient()
	res, err := mg.Update("user", u)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestMgoClient_Delete(t *testing.T) {
	u := &User{Username: "testUser", Nickname: "testNickname", ID: uid}
	mg := getMgoClient()
	res, err := mg.Delete("user", u)
	assert.Nil(t, err)
	assert.True(t, res)
}
