package mongodb

import (
	"encoding/json"
	"github.com/luo-zn/db-adaptor/bases"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//BsonMarshal convert an instance to []byte.
//It return []byte, error
func BsonMarshal(e interface{}) ([]byte, error) {
	mp, er := bases.Entity2Map(e)
	if er != nil {
		return nil, er
	}
	if ids, ok := mp["id"].(string); ok && ids != "" {
		id, er1 := primitive.ObjectIDFromHex(ids)
		if er1 != nil {
			return nil, er1
		}
		delete(mp, "id")
		mp["_id"] = id
	}
	return bson.Marshal(mp)
}

//Map2ClientOptions convert json to ClientOptions of mongo-driver.
// It return *options.ClientOptions, error
func Map2ClientOptions(opt map[string]interface{}) (*options.ClientOptions, error) {
	optByte, er1 := json.Marshal(opt)
	if er1 != nil {
		return nil, er1
	}
	op := options.Client()
	er2 := json.Unmarshal(optByte, &op)
	if er2 != nil {
		return nil, er2
	}
	uri := op.GetURI()
	if uri == "" {
		constr, ok := opt["uri"].(string)
		if !ok {
			return nil, bases.Error("Mongo connecting needs uri string!")
		}
		uri = constr
	}
	return op.ApplyURI(uri), nil
}
