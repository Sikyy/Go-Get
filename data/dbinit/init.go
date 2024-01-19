package dbinit

import (
	"Go-Get/data"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

var Client *mongo.Client

func init() {
	// 建立MongoDB连接
	client, err := data.ConnectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	// 设置全局变量
	Client = client

	fmt.Println("MongoDB连接成功")
}

// 在应用程序退出时执行断开连接操作
func CloseMongoDBConnection() {
	if Client != nil {
		err := Client.Disconnect(context.Background())
		if err != nil {
			log.Println("Error disconnecting from MongoDB:", err)
		} else {
			fmt.Println("Disconnected from MongoDB")
		}
	}
}
