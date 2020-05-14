package bases

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Recover() {
	if r := recover(); r != nil {
		if err, ok := r.(error); !ok {
			fmt.Print(err)
		}
	}
}

func Entity2Map(e interface{}) (map[string]interface{}, error) {
	f, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	var mp map[string]interface{}
	if err1 := json.Unmarshal(f, &mp); err1 != nil {
		return nil, err1
	}
	return mp, nil
}

func CheckMongoID(e interface{}) (map[string]interface{}, error) {
	mp, _ := e.(map[string]interface{})
	ids, ok := mp["id"].(string)
	if !ok {
		return mp, Error("ID is not a string.")
	}
	if ids != "" {
		id, er1 := primitive.ObjectIDFromHex(ids)
		if er1 != nil {
			return mp, er1
		}
		mp["_id"] = id
	}
	delete(mp, "id")
	return mp, nil
}