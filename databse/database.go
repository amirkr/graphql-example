package database


import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"context"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
	"gitlab.com/amirkerroumi/my-gqlgen/graph/model"
)
type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
        log.Fatal("Mongo NewClient error:", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Failure to Connect to MongoDB:", err.Error())
	}

	return &DB {
		client: client,
	}
}

func (db* DB) Save(input *model.NewAuthor) *model.Author {
	collection := db.client.Database("my-gqlgen").Collection("author")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Fatal("MongoDB Author Insertion Failure:", err.Error())
	}
	return &model.Author {
		ID: res.InsertedID.(primitive.ObjectID).Hex(),
		Firstname: input.Firstname,
		Lastname: input.Lastname,
	}
}

func (db *DB) FindyID(ID string) *model.Author {
	ObjectID, err := primitive.ObjectIDFromHex(ID)
	collection := db.client.Database("my-gqlgen").Collection("author")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := collection.Find(ctx, bson.M{"_id": ObjectID})
	if err != nil {
		log.Fatal("MongoDB Author Find Failure:", err.Error())
	}
	var author *model.Author
	res.Decode(&author)
	return author
}

func (db *DB) All() []*model.Author {
	var authors []*model.Author
	collection := db.client.Database("my-gqlgen").Collection("author")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal("MongoDB Author Find Failure:", err.Error())
	}
	for cursor.Next(ctx){
		var author *model.Author
		err := cursor.Decode(&author)
		if err != nil {
			log.Fatal("MongoDB Failure to Decode Author:", err.Error())
		}
		authors = append(authors, author)
	}

	return authors
}