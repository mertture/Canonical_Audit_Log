package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mertture/audit-log/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


func (server *Server) CreateEvent(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	event := models.Event{}

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set timestamp if not provided
	if event.EventTime.IsZero() {
		event.EventTime = time.Now()
	}

	// Insert event into MongoDB
	collection := server.DB.Collection("Event")
	result, err := collection.InsertOne(ctx, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID, "message": "Successfully created the event"})
}

func (server *Server) GetAllEvents(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Query events from MongoDB
	collection := server.DB.Collection("Event")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	// Decode all events into a slice
    var events []models.Event
    if err := cursor.All(ctx, &events); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while decoding the events"})
        return
    }

	// Return list of events
	c.JSON(http.StatusOK, events)
}

func (server *Server) GetEventByID(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Get event ID from URL parameter
	eventID := c.Param("id")

	// Query event from MongoDB by ID
	collection := server.DB.Collection("Event")
	filter := bson.M{"_id": eventID}

	event := models.Event{}
	err := collection.FindOne(ctx, filter).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Return event as JSON response
	c.JSON(http.StatusOK, event)
}

func (server *Server) DeleteEvent(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Get event ID from URL parameter
	eventID := c.Param("id")

	// Delete event from MongoDB by ID
	collection := server.DB.Collection("Event")
	filter := bson.M{"_id": eventID}
	
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if event was deleted
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted the event"})
}

