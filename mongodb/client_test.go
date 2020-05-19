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
	"log"
	"reflect"
	"testing"
	"time"
)

var (
	dbUri = "mongodb://localhost:27017"
	uid   = "5eb8e1f6f715a2494825d0ba"
)

type User struct {
	ID       string            `json:"id,omitempty" bson:"_id,omitempty"`
	Username string            `json:"username,omitempty" bson:"username,omitempty"`
	Nickname string            `json:"nickname,omitempty" bson:"nickname,omitempty"`
	Region   map[string]string `json:"region,omitempty" bson:"region,omitempty"`
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

func TestMgoClient(t *testing.T) {
	var mgc *mongo.Client
	var mgcoll *mongo.Collection
	guardCon := monkey.Patch(mongo.Connect,
		func(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
			c, err := mongo.NewClient(opts...)
			if err != nil {
				return nil, err
			}
			log.Print("calling monkey.Connect")
			return c, nil
		})
	mg := getMgoClient()
	u := &User{Username: "testUser", Nickname: "testNickname", ID: uid}
	t.Run("Connect", func(t *testing.T) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(mgc), "Ping",
			func(_ *mongo.Client, ctx context.Context, rp *readpref.ReadPref) error {
				return nil
			})
		defer guard.Unpatch()
		err := mg.client.Ping(context.TODO(), readpref.SecondaryPreferred())
		assert.Nil(t, err)
		assert.NotEmpty(t, mg.client)
	})
	t.Run("getCollection", func(t *testing.T) {
		col := mg.getCollection("testDB", "user")
		assert.NotNil(t, col)
		assert.Equal(t, col.Name(), "user")
	})
	t.Run("Create", func(t *testing.T) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(mgcoll), "InsertOne",
			func(_ *mongo.Collection, ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
				return &mongo.InsertOneResult{InsertedID: uid}, nil
			})
		defer guard.Unpatch()
		res, err := mg.Create("user", u)
		assert.NotNil(t, res)
		assert.Nil(t, err)
	})
	t.Run("FindOne", func(t *testing.T) {
		var mgc *mongo.SingleResult
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(mgc), "Decode",
			func(_ *mongo.SingleResult, v interface{}) error {
				obj, ok := v.(*User)
				if ok {
					obj.ID = uid
					return nil
				}
				return bases.Error("Obj is not User struct!")
			})
		defer guard.Unpatch()
		_, err := mg.FindOne("user", u)
		assert.Nil(t, err, "Error %s", err)
		assert.Equal(t, "testUser", u.Username, "The expected Username is testUser")
		assert.Equal(t, "testNickname", u.Nickname, "The expected Nickname is testNickname")
		assert.NotEmpty(t, u.ID, "User id is %s", u.ID)
	})
	t.Run("Retrieve", func(t *testing.T) {
		var mgc *mongo.SingleResult
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(mgc), "Decode",
			func(_ *mongo.SingleResult, v interface{}) error {
				obj, ok := v.(*User)
				if ok {
					obj.ID = uid
					return nil
				}
				return bases.Error("Obj is not User struct!")
			})
		defer guard.Unpatch()
		u.ID = ""
		_, err := mg.Retrieve("user", u)
		assert.Nil(t, err, "Error %s", err)
		assert.Equal(t, u.Username, "testUser", "The expected Username is testUser")
		assert.Equal(t, u.Nickname, "testNickname", "The expected Nickname is testNickname")
		assert.NotEmpty(t, u.ID, "User id is %s", u.ID)
		assert.Equal(t, uid, u.ID, "The expected ID is %s", uid)
	})
	t.Run("UpdateOneById", func(t *testing.T) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(mgcoll), "UpdateOne",
			func(_ *mongo.Collection, ctx context.Context, filter interface{}, update interface{},
				opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
				return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
			})
		defer guard.Unpatch()
		res, err := mg.UpdateOneById("user", uid, u)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("Update", func(t *testing.T) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(mgcoll), "UpdateOne",
			func(_ *mongo.Collection, ctx context.Context, filter interface{}, update interface{},
				opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
				return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
			})
		defer guard.Unpatch()
		upd := u
		res, err := mg.Update("user", u, upd)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("UpdateOneWithFilter", func(t *testing.T) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(mgcoll), "UpdateOne",
			func(_ *mongo.Collection, ctx context.Context, filter interface{}, update interface{},
				opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
				return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
			})
		defer guard.Unpatch()
		filter := map[string]interface{}{"username": u.Username, "nickname": u.Nickname}
		u.Region = map[string]string{"country": "China", "city": "shenzhen"}
		res, err := mg.UpdateOneWithFilter("user", filter, u)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("Count", func(t *testing.T) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(mgcoll), "CountDocuments",
			func(_ *mongo.Collection, ctx context.Context, filter interface{},
				opts ...*options.CountOptions) (int64, error) {
				return 1, nil
			})
		defer guard.Unpatch()
		res, err := mg.Count("user", u)
		assert.NotNil(t, res)
		assert.Greaterf(t, res, int64(0), "")
		assert.Nil(t, err)
	})
	t.Run("Delete", func(t *testing.T) {
		guardDelOne := monkey.PatchInstanceMethod(reflect.TypeOf(mgcoll), "DeleteOne",
			func(_ *mongo.Collection, ctx context.Context, filter interface{},
				opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
				log.Print("calling monkey.DeleteOne", filter)
				return &mongo.DeleteResult{DeletedCount: 1}, nil
			})
		defer guardDelOne.Unpatch()
		res, err := mg.Delete("user", u)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	defer guardCon.Unpatch()
	//mg.Close()
}
