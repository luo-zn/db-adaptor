package mongodb

import (
	"context"
	"db-adaptor/bases"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (m *MgoClient) connect(uri string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // bug may happen
	if client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri)); err == nil {
		m.ctx = ctx
		m.client = client
	}
}

func (m *MgoClient) Connect(uri string, timeout time.Duration) {
	if timeout == 0 {
		timeout = 20 * time.Second
	}
	logrus.Info("Connect count!!!!!!!")
	m.connect(uri, timeout)
}

func (m *MgoClient) close() error {
	return m.client.Disconnect(m.ctx)
}

func (m *MgoClient) getCollection(db string, tb string) *mongo.Collection {
	return m.client.Database(db).Collection(tb)
}

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

func (m *MgoClient) DeleteOne(tb string, e bases.Entity) (bool, error) {
	defer bases.Recover()
	f, er1 := bson.Marshal(e)
	if er1 != nil {
		return false, er1
	}
	res, err := m.getCollection(e.DataBase(), tb).DeleteOne(context.TODO(), f)
	if err != nil {
		return false, err
	}
	return res.DeletedCount == 1, nil
}

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

func (m *MgoClient) Retrieve(tb string, filter bases.Entity) (interface{}, error) {
	return m.FindOne(tb, filter)
}

func (m *MgoClient) Update(tb string, e bases.Entity) (bool, error) {
	defer bases.Recover()
	mp, er := bases.Entity2Map(e)
	if er != nil {
		return false, er
	}
	ids, ok := mp["id"].(string)
	if !ok {
		return false, bases.Error("ID is not a string.")
	}
	//id, er1 := primitive.ObjectIDFromHex(ids)
	//if er1 != nil {
	//	return false, er1
	//}
	filter := bson.M{"_id": ids}
	delete(mp, "id")
	update, _ := bson.Marshal(bson.M{"$set": mp})
	res, err := m.getCollection(e.DataBase(), tb).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount == 1, nil
}

func (m *MgoClient) Delete(tb string, e bases.Entity) (bool, error) {
	return m.DeleteOne(tb, e)
}
