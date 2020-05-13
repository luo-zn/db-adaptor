package mongodb

//SetWrapper Used for mongodb update with $set
type SetWrapper struct {
	Set interface{} `bson:"$set,omitempty"`
}
