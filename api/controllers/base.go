package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	DB  *mongo.Database
	Router  *gin.Engine
}

func (server *Server) Initialize(DBurl string) {
	var err error 

	//DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	// Set up MongoDB client options
	clientOptions := options.Client().ApplyURI(DBurl)

	// Connect to MongoDB client
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Printf("Cannot connect to database")
		log.Fatal("This is the error:", err)
	}

	// Check the MongoDB connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Printf("Cannot connect to database")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the database")
	}

	// Set up MongoDB database and collections
	server.DB = client.Database("audit-log")

	server.Router = gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")

	server.Router.Use(cors.New(config))

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}