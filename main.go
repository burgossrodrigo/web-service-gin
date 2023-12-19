package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type user struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"` // Password should be a string
}

var collection *mongo.Collection
var userCollection *mongo.Collection

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Replace <password> with your actual password and update the URI as needed
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(""))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection = client.Database("golang").Collection("albums")
	userCollection = client.Database("golang").Collection("users")

	router := gin.Default()

	router.Static("/static", "./static")

	router.POST("/users", postUser)
	router.POST("/login", loginUser)

	router.GET("/create-account", func(c *gin.Context) {
		c.File("./static/CreateAccount.html") // Update the path according to your file structure
	})

	router.GET("/logged-in", func(c *gin.Context) {
		c.File("./static/Home.html")
	})

	router.GET("/signin", func(c *gin.Context) {
		c.File("./static/Login")
	})

	router.Run("localhost:8080")
}

func postUser(c *gin.Context) {
	// Extract data from form submission
	email := c.PostForm("email")
	name := c.PostForm("name")
	password := c.PostForm("password")

	// TODO: Validate the inputs and hash the password

	newUser := user{
		Email:    email,
		Name:     name,
		Password: password, // Remember to hash the password before storing
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// You may want to redact the password before sending the response
	newUser.Password = ""
	c.IndentedJSON(http.StatusCreated, newUser)
}

func loginUser(c *gin.Context) {
	// Extract login data from the form
	email := c.PostForm("email")
	password := c.PostForm("password") // In a real-world app, compare hashed passwords

	var user user
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Search for user by email
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// TODO: Verify the password, compare with hashed password stored in DB

	// For demonstration, assuming passwords match if they are exactly the same
	if user.Password != password {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Redirect to a logged-in page
	c.Redirect(http.StatusFound, "/logged-in")
}
