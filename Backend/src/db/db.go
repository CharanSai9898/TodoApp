package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database

func ConnectDB(){
	
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = os.Getenv("DB_Name")
	}

	if mongoURI ==  "" || dbName == ""{
		log.Fatal("DB Credentials are missing")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	client,err := mongo.Connect(ctx,clientOptions)
	if err != nil{
		log.Fatal("Failed to connect to MongoDb:",err)
	}

	err = client.Ping(ctx,nil)
	if err != nil{
		log.Fatal("MongoDB ping Failed",err)
	}
	fmt.Println("Connected to MongoDB")
	Db= client.Database(dbName)
}
func GetCollection(name string) *mongo.Collection{
	return Db.Collection(name)
}