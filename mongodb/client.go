package mongodb

import (
	"context"
	"github.com/luo-zn/db-adaptor/bases"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func NewMgoClient(uri string) *MgoClient {
	client, _ := mongo.NewClient(options.Client().ApplyURI(uri))
	return &MgoClient{client: client}
}

type MgoClient struct {
	client *mongo.Client
	ctx    context.Context
}

//connect call mongo.Connect.
func (m *MgoClient) connect(opt map[string]interface{}) error {
	clientOpt, er1 := Map2ClientOptions(opt)
	if er1 != nil {
		return er1
	}
	ctxTimeout, ok := opt["ctx_timeout"].(time.Duration)
	if !ok {
		ctxTimeout = 40 * time.Second
	}
	socketTimeout, ok := opt["sockettimeout"].(time.Duration)
	if !ok {
		socketTimeout = 30 * time.Minute
	}
	clientOpt.SocketTimeout = &socketTimeout
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel() // bug may happen
	if client, err := mongo.Connect(ctx, clientOpt); err == nil {
		m.ctx = ctx
		m.client = client
	}
	return nil
}

//Connect call MgoClient.connect
func (m *MgoClient) Connect(opt map[string]interface{}) error {
	return m.connect(opt)
}

//Close call MgoClient。client.Disconnect
func (m *MgoClient) Close() error {
	return m.client.Disconnect(m.ctx)
}

//getCollection call mongo.Collection.
func (m *MgoClient) getCollection(db string, tb string) *mongo.Collection {
	return m.client.Database(db).Collection(tb)
}

//FindOne call mongo.Collection.FindOne
func (m *MgoClient) FindOne(tb string, e bases.Entity) (bases.Entity, error) {
	defer bases.Recover()
	filter, er1 := bson.Marshal(e)
	//filter,er1 := BsonMarshal(e)
	if er1 != nil {
		return nil, er1
	}
	if err := m.getCollection(e.DataBase(), tb).FindOne(context.TODO(), filter).Decode(e); err != nil {
		return nil, err
	}
	return e, nil
}

//UpdateOneById call mongo.Collection.UpdateOne.
func (m *MgoClient) UpdateOneById(tb string, ids string, u bases.Entity) (bool, error) {
	defer bases.Recover()
	update, er1 := bson.Marshal(SetWrapper{Set: u})
	if er1 != nil {
		return false, er1
	}
	//id, er2 := primitive.ObjectIDFromHex(ids)
	//if er2 != nil{
	//	return false, er2
	//}
	res, err := m.getCollection(u.DataBase(), tb).UpdateOne(context.TODO(), bson.M{"_id": ids}, update)
	if err != nil {
		return false, err
	}
	return res.MatchedCount == 1 || res.ModifiedCount == 1, nil
}

//DeleteOne call mongo.Collection.DeleteOne.
//The e parameter is type of bases.Entity.
func (m *MgoClient) DeleteOne(tb string, e bases.Entity) (bool, error) {
	defer bases.Recover()
	f, er1 := bson.Marshal(e)
	if er1 != nil {
		return false, er1
	}
	coll := m.getCollection(e.DataBase(), tb)
	log.Print("coll=", coll)
	res, err := coll.DeleteOne(context.TODO(), f)
	log.Print("DeleteOne=", res, err)
	if err != nil {
		return false, err
	}
	return res.DeletedCount == 1, nil
}

//Create call mongo.Collection.InsertOne.
//The e parameter is type of bases.Entity.
func (m *MgoClient) Create(tb string, e bases.Entity) (interface{}, error) {
	var err error
	defer bases.Recover()
	collection := m.getCollection(e.DataBase(), tb)
	res, err := collection.InsertOne(context.TODO(), e)
	if err == nil {
		id, ok := res.InsertedID.(primitive.ObjectID)
		if ok {
			return id.Hex(), nil
		} else {
			return res.InsertedID, nil
		}
	}
	return "", err
}

//Retrieve call .MgoClient.FindOne
func (m *MgoClient) Retrieve(tb string, filter bases.Entity) (interface{}, error) {
	return m.FindOne(tb, filter)
}

//Update call mongo.Collection.UpdateOne.
//The filter parameter is type of bases.Entity, update parameter is type of  bases.Entity
func (m *MgoClient) Update(tb string, f bases.Entity, e bases.Entity) (bool, error) {
	defer bases.Recover()
	//mp, er := bases.Entity2Map(e)
	//if er != nil {
	//	return false, er
	//}
	//ids, ok := mp["id"].(string)
	//if !ok {
	//	return false, bases.Error("ID is not a string.")
	//}
	//id, er1 := primitive.ObjectIDFromHex(ids)
	//if er1 != nil {
	//	return false, er1
	//}
	//filter := bson.M{"_id": ids}
	//delete(mp, "id")
	//update, _ := bson.Marshal(bson.M{"$set": mp})
	update, er1 := bson.Marshal(SetWrapper{Set: e})
	if er1 != nil {
		return false, er1
	}
	filter, er2 := bson.Marshal(f)
	if er2 != nil {
		return false, er2
	}
	res, err := m.getCollection(e.DataBase(), tb).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount == 1, nil
}

//UpdateOneWithFilter call mongo.Collection.UpdateOne.
//filter parameter is type of map[string]interface{}, update parameter is type of  bases.Entity
func (m *MgoClient) UpdateOneWithFilter(tb string, filter map[string]interface{}, e bases.Entity) (bool, error) {
	defer bases.Recover()
	//mp, er := bases.Entity2Map(e)
	//if er != nil {
	//	return false, er
	//}
	//delete(mp, "_id")
	//update, _ := bson.Marshal(bson.M{"$set": mp})
	update, er1 := bson.Marshal(SetWrapper{Set: e})
	if er1 != nil {
		return false, er1
	}
	res, err := m.getCollection(e.DataBase(), tb).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount == 1, nil
}

//Count call mongo.Collection.CountDocuments to count documents, bases.Entity instance as filter parameter.
func (m *MgoClient) Count(tb string, f bases.Entity) (int64, error) {
	defer bases.Recover()
	filter, er1 := bson.Marshal(f)
	if er1 != nil {
		return 0, er1
	}
	res, err := m.getCollection(f.DataBase(), tb).CountDocuments(context.TODO(), filter)
	return res, err
}

//Delete call MgoClient.DeleteOne.
func (m *MgoClient) Delete(tb string, e bases.Entity) (bool, error) {
	return m.DeleteOne(tb, e)
}
