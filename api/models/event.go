package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
    ID            primitive.ObjectID     `json:"id" bson:"_id"`
    EventType     int                    `json:"event_type" bson:"event_type"`
    EventTime     time.Time              `json:"event_time" bson:"event_time"`
    UserID        primitive.ObjectID     `json:"user_id" bson:"user_id"`
    ServiceName   string                 `json:"service_name" bson:"service_name"`
    Status        int                    `json:"status" bson:"status"`
    EventFields   map[string]interface{} `json:"event_fields" bson:"event_fields"`
}

type EventPayload struct {
    EventType     string                 `json:"event_type" bson:"event_type"`
    EventTime     time.Time              `json:"event_time" bson:"event_time"`
    UserID        primitive.ObjectID     `json:"user_id" bson:"user_id"`
    ServiceName   string                 `json:"service_name" bson:"service_name"`
    Status        string                 `json:"status" bson:"status"`
    EventFields   map[string]interface{} `json:"event_fields" bson:"event_fields"`
}

func (e *Event) Prepare() {
	e.ID = primitive.NewObjectID();
}

func (payload *EventPayload) ToEvent(eventType int, status int) Event {
    event := Event{
        EventType: eventType,
        EventTime: payload.EventTime,
		UserID: payload.UserID,
        ServiceName: payload.ServiceName,
		Status: status,
        EventFields: payload.EventFields,
    }
    return event
}