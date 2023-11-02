package dbinit

import (
	"Go-Get/data"
	"context"
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
	defer client.Disconnect(context.Background())
}
