package mongodbpkg

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var MongoDBCon mongoCon

type mongoCon struct {
	DB *mongo.Client
}

func NewMongoGoDriverCon(url, dbName string) error {

	db, err := dialMongodb(url, dbName)
	if err != nil {
		return err
	}

	MongoDBCon.DB = db

	return nil
}

func dialMongodb(url, dbName string) (*mongo.Client, error) {

	ops := options.Client().ApplyURI(url)
	ops.SetReadPreference(readpref.SecondaryPreferred())
	ops.SetMaxPoolSize(5000)
	ops.SetConnectTimeout(5 * time.Second)
	ops.SetMaxConnIdleTime(30 * time.Second)

	client, err := mongo.NewClient(ops)
	if err != nil {
		return nil, err
	}
	err = client.Connect(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("can't connect to database %v %v", dbName, err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't ping to database %v %v", dbName, err)
	}
	return client, nil
}
