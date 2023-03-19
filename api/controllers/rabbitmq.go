package controllers

import (
	"encoding/json"
	"context"
	"time"
	"fmt"
	"log"
	"github.com/mertture/audit-log/api/models"
	"github.com/streadway/amqp"

)

// Queue represents a connection to a RabbitMQ server and a channel for publishing messages
type Queue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// Init initializes a new Queue struct with a connection to a RabbitMQ server
func Init(uri string, queueName string) (*Queue, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare RabbitMQ queue: %v", err)
	}

	return &Queue{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (q *Queue) Push(event models.Event) error {
    // Serialize event data
    eventData, err := json.Marshal(event)
    if err != nil {
        return err
    }

    // Push event data to the queue
    err = q.channel.Publish(
        "",         // exchange
        q.queue.Name, // routing key
        false,      // mandatory
        false,      // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        eventData,
        },
    )
    if err != nil {
        return err
    }

    return nil
}

func (server *Server) Consume() error {
    msgs, err := server.Queue.channel.Consume(
        server.Queue.queue.Name, // queue
        "",          // consumer
        false,       // auto-ack
        false,       // exclusive
        false,       // no-local
        false,       // no-wait
        nil,         // args
    )
    if err != nil {
        return err
    }

    for msg := range msgs {
        // Deserialize event data
        var event models.Event
        err := json.Unmarshal(msg.Body, &event)
        if err != nil {
            log.Println("Error deserializing message:", err)
            msg.Ack(false)
            continue
        }
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

        // Insert event into MongoDB
        collection := server.DB.Collection("Event")
        _, err = collection.InsertOne(ctx, event)
        if err != nil {
            log.Println("Error inserting event into MongoDB:", err)
            msg.Ack(false)
            continue
        }

        // Acknowledge message
        msg.Ack(false)
    }

    return nil
}


// Close gracefully closes the connection to the RabbitMQ server
func (q *Queue) Close() error {
	err := q.channel.Close()
	if err != nil {
		return fmt.Errorf("failed to close RabbitMQ channel: %v", err)
	}

	err = q.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close RabbitMQ connection: %v", err)
	}

	return nil
}
