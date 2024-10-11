package mongo

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
	URI    string
}

func NewMongoDBClient(uri string) *MongoClient {
	return &MongoClient{
		URI: uri,
	}
}

func (m *MongoClient) Connect() error {
	err := mgm.SetDefaultConfig(nil, "app", options.Client().ApplyURI(m.URI))
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoClient) Disconnect() error {
	if m.Client != nil {
		return m.Client.Disconnect(nil)
	}
	return nil
}
