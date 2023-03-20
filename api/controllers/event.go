package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mertture/audit-log/api/constants"
	"github.com/mertture/audit-log/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


func (server *Server) CreateEvent(c *gin.Context) {

    payload := models.EventPayload{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set timestamp if not provided
	if payload.EventTime.IsZero() {
		payload.EventTime = time.Now()
	}


	// Convert event type string to number
	eventType, ok := constants.EventTypeNumbers[payload.EventType]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event type"})
		return
	}

	// Convert status string to number
	status, ok := constants.StatusNumbers[payload.Status]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}


	// Convert payload to event
	event := payload.ToEvent(eventType, status)

	event.Prepare()

	err := server.Queue.Push(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully added create event to the MQ"})
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
	eventID, err := primitive.ObjectIDFromHex(c.Param("id"))

	// Query event from MongoDB by ID
	collection := server.DB.Collection("Event")
	filter := bson.M{"_id": eventID}
	event := models.Event{}
	err = collection.FindOne(ctx, filter).Decode(&event)
	fmt.Println("evv:", event);
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

// Bulk get event type
func (server *Server) GetEventByTypeID(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Get event ID from URL parameter
	eventTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
    	c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot convert event type to int"})
		return
	}

	// Query event from MongoDB by ID
	collection := server.DB.Collection("Event")
	filter := bson.M{"event_type": eventTypeID}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	defer cursor.Close(ctx)

	// Decode all type events into a slice
    var events []models.Event
    if err := cursor.All(ctx, &events); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while decoding the events"})
        return
    }

	// Return list of events
	c.JSON(http.StatusOK, events)
}


func (server *Server) DeleteEvent(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Get event ID from URL parameter
	eventID, err := primitive.ObjectIDFromHex(c.Param("id"))


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

