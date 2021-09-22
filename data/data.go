package data

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var err error

func init() {
	password := os.Getenv("PASS")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://pass:" + password + "@cluster0.zlr5c.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Println("Error connecting to client:", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Printf("Error connecting to mongo server: %v", err)
	}

	collection = client.Database("users").Collection("user")
}

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Password string             `bson:"password"`
}

func (u *User) Store() bool {
	var result bson.M
	err := collection.FindOne(context.Background(), bson.D{{"name", u.Name}}).Decode(&result)
	if result["name"].(string) == u.Name {
		return false
	}
	_, err = collection.InsertOne(context.Background(), u)
	if err != nil {
		log.Printf("Error inserting person record: %v", err)
		return false
	}
	return true
}

func (u *User) Auth() bool {
	var result bson.M
	err := collection.FindOne(context.Background(), bson.D{{"name", u.Name}}).Decode(&result)
	if err != nil {
		log.Println("User not authorized")
		// return false
	}
	if result["name"].(string) == u.Name && result["password"].(string) == u.Password {
		return true
	} else {
		return false
	}
}
