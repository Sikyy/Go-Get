package data

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectToMongoDB() (*mongo.Client, error) {
	url := "mongodb://localhost:27017"

	// 设置客户端连接配置
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	// 检查连接
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("Could not connect to MongoDB:", err)
		return nil, err
	}

	return client, nil
}

func InsertDocument(client *mongo.Client, uploadinfo bson.M, datevase string, collections string) error {
	// 获取要插入的数据库和集合
	collection := client.Database(datevase).Collection(collections)
	if collection == nil {
		return errors.New("collection is nil")
	}

	// 插入文档
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, uploadinfo)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Document inserted successfully")
	return nil
}

func QueryDocuments(client *mongo.Client) {
	// 获取要查询的数据库和集合
	collection := client.Database("Go-Get-MongoDB").Collection("test")

	// 构建查询条件
	filter := bson.M{
		"name": "John Doe",
	}

	// 查询文档
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// 遍历查询结果
	for cursor.Next(ctx) {
		var result map[string]interface{}
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		log.Printf("Found document: %v", result)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}
