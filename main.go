package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Title    string             `json:"title" bson:"title"`
	Author   Author  `json:"author" bson:"author"`
	Quantity int                `json:"quantity" bson:"quantity"`
}

type Author struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

var client *mongo.Client

func main() {
	r := gin.Default() // router variable

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	r.GET("/books", getAllBooks)
	// r.GET("/users/:id", getSingleUser)
	r.POST("/books", createBook)

	r.Run(":8080")
}

func ConnectToMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") // Change the URI to your MongoDB server URI
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getAllBooks(c *gin.Context) {
	collection := client.Database("testdb").Collection("books") // Replace "yourdbname" with your database name
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	var books []Book
	for cur.Next(context.Background()) {
		var book Book
		if err := cur.Decode(&book); err != nil {
			log.Printf("Error decoding book: %v", err)
			continue
		}
		books = append(books, book)
	}
	c.JSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook Book

	// Bind the JSON data from the request body to the newBook struct
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBook.ID = primitive.NewObjectID()
	newBook.Author.ID = primitive.NewObjectID()

	// Get the "books" collection
	collection := client.Database("testdb").Collection("books") // Replace "testdb" with your actual database name

	// Insert the new book into the collection without specifying the _id
	_, err := collection.InsertOne(context.TODO(), newBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the actual data of the newly created book
	c.JSON(http.StatusCreated, newBook)
}
