package main

import (
	"context"
	"flag"
	"github.com/globalsign/mgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	url = "mongodb://my_user:password123@localhost:27018/orderbook"
)

func main() {
	count := flag.Int("count", 10000, "count")
	flag.Parse()
	data := make([]interface{}, 0, *count)
	for i := 0; i < *count; i++ {
		data = append(data, TestData{})
	}
	if err := official("test", "test", data); err != nil {
		panic(err.Error())
	}

	if err := clean("test", "test"); err != nil {
		panic(err.Error())
	}

	if err := community("test", "test", data); err != nil {
		panic(err.Error())
	}
	if err := clean("test", "test"); err != nil {
		panic(err.Error())
	}
	if err := officialBatch("test", "test", data); err != nil {
		panic(err.Error())
	}
	if err := clean("test", "test"); err != nil {
		panic(err.Error())
	}
	if err := communityBatch("test", "test", data); err != nil {
		panic(err.Error())
	}
	if err := clean("test", "test"); err != nil {
		panic(err.Error())
	}
}

func official(db, collectionName string, data []interface{}) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	collection := client.Database(db).Collection(collectionName)
	start := time.Now()
	for _, element := range data {
		if _, err = collection.InsertOne(ctx, element); err != nil {
			return err
		}

	}
	println("官方包耗时" + time.Since(start).String())

	return nil
}
func officialBatch(db, collectionName string, data []interface{}) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	collection := client.Database(db).Collection(collectionName)
	start := time.Now()

	if _, err = collection.InsertMany(ctx, data); err != nil {
		return err
	}

	println("官方包批量耗时" + time.Since(start).String())

	return nil
}
func community(db, collectionName string, data []interface{}) error {
	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}

	collection := session.DB(db).C(collectionName)
	start := time.Now()

	for _, element := range data {
		if err = collection.Insert(element); err != nil {
			return err
		}
	}
	println("社区包耗时" + time.Since(start).String())

	return nil
}

func communityBatch(db, collectionName string, data []interface{}) error {
	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}

	collection := session.DB(db).C(collectionName)
	start := time.Now()

	if err = collection.Insert(data...); err != nil {
		return err
	}

	println("社区包批量耗时" + time.Since(start).String())

	return nil
}
func clean(db, collectionName string) error {
	session, err := mgo.Dial(url)

	if err != nil {
		return err
	}

	return session.DB(db).C(collectionName).DropCollection()

}

type TestData struct {
	A int
	B string
	C bool
}
