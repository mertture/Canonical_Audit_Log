package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
    ID            primitive.ObjectID     `json:"id"`
    EventType     string                 `json:"event_type"`
    EventTime     time.Time              `json:"event_time"`
    CommonFields  map[string]interface{} `json:"common_fields"`
    EventFields   map[string]interface{} `json:"event_fields"`
}

type CommonFields struct {
    UserID    primitive.ObjectID `json:"user_id"`
}

type CustomerCreatedEvent struct {
    CommonFields
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
}

type CustomerActionPerformedEvent struct {
    CommonFields
    ActionName string `json:"action_name"`
    ResourceID string `json:"resource_id"`
}

type CustomerBilledEvent struct {
    CommonFields
    Amount float64 `json:"amount"`
}

type CustomerDeactivatedEvent struct {
    CommonFields
    Reason string `json:"reason"`
}