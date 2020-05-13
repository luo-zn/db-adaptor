/* Create By Jenner.luo */
package mongodb

import (
	"db-adaptor/bases"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
