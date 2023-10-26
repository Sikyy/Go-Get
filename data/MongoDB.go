package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type FileUpload struct {
	Name        string
	Size        int
	FileCount   int
	CreatedDate time.Time
	UploadDate  time.Time
	Uploader    string
	Hash        string
	Tracker     string
	Comment     string
}

func ConnectToMongoDB() (*mongo.Client, error) {
	url := "mongodb://localhost:27017"
	if url == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable")
	}
	// 设置客户端连接配置
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	//在程序结束时断开与MongoDB的连接
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// 检查连接
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	}

	// 连接成功输出
	fmt.Println("Successfully connected to MongoDB")

	// 返回连接的客户端
	return client, nil
}

func InsertFileToDatabase(client *mongo.Client, filePath string) error {

	// 选择数据库和集合
	database := client.Database("Go-Get-MongoDB")
	collection := database.Collection("Torrents")

	// 创建一个 FileUpload 实例
	upload := FileUpload{
		Name:        "example.txt",           //文件名
		Size:        1024,                    //文件大小
		FileCount:   1,                       //文件数量
		CreatedDate: time.Now(),              //创建时间
		UploadDate:  time.Now(),              //上传时间
		Uploader:    "user123",               //上传者
		Hash:        "abc123",                //哈希值
		Tracker:     "http://tracker.com",    //Tracker地址
		Comment:     "This is a test upload", //备注
	}

	// 插入文档
	_, err := collection.InsertOne(context.TODO(), upload)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Document inserted successfully")

	return nil
}
