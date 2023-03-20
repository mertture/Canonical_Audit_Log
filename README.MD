# Audit Logging Microservice API

This microservice API is designed to receive, store, and retrieve audit log events. It is built with Golang and MongoDB, and uses RabbitMQ for write-intensive operations.

## How to Run for the first time

**Install and Start RabbitMQ**

## For MacOS

``` bash
$ brew install rabbitmq
```
Running RabbitMQ on local
``` bash
$ brew services start rabbitmq
```

Stopping RabbitMQ on local
``` bash
$ brew services stop rabbitmq
```

## For Ubuntu 

You can follow the documentation

https://www.rabbitmq.com/install-debian.html

## Continuing with cloning the repository

**Clone the repository**
```bash
git clone https://github.com/mertture/canonical_audit_log.git
```

**Install the go packages**
``` bash
$ cd location-of-repo
$ go mod download
```

**Run the Service**
```bash
$ go run main.go
```

### Single Command For Starting After Running the RabbitMQ Server
```bash
git clone https://github.com/mertture/canonical_audit_log.git
cd canonical_audit_log
go mod download
go run main.go
```


## Config Variables

```bash
DB_URL=<your-mongodb-uri>
AMQP_URL=<your-rabbitmq-uri>
```

## Features

- User authentication (JWT) using email and password
- Creation and retrieval of audit log events
- Retrieval of audit log events by event type
- Event payload validation using JSON schema
- Architecture

The microservice is designed using the following architecture:

- **API layer**: Handles HTTP requests and responses.
- **Service layer**: Contains business logic and validation.
- **Message Queue layer**: Communicates with RabbitMQ.

**Architecture-Diagram**

## API Endpoints

- POST /api/users/register
- POST /api/users/login
- POST /api/events
- GET /api/events
- GET /api/events/:id
- DELETE /api/events/:id
- GET /api/events/types/:id


### **User Endpoints**


### **User Registration**


- POST /api/users/register

Registers a new user with the provided email and password.

**Request Payload**
```
{
    "email": <email>,
    "password": <string>
}
```
**Response Payload**
```
{
    "id": <string>,
}
```

### **User Login Endpoint**

- POST /api/users/login

This endpoint is used for authenticating registered users. It accepts a JSON payload containing the user's email and password, and returns a JSON Web Token (JWT) if the email and password are valid.

**Request Body**
```
{
    "email": <email>,
    "password": <string>
}
```
**Response:**

   - **Status code:** 200 OK

   - **Body:**

```
{
    "token": <string>
}
```


**Possible Errors:**

   - **Status code:** 400 Bad Request

   - **Body:**

```
{
    "message": "Invalid request payload"
}
```



   - **Status code:** 401 Unauthorized

   - **Body:**
```
{
    "message": "Invalid email or password"
}
```



   - **Status code:** 500 Internal Server Error

   - **Body:**
```
{
    "message": "Internal server error"
}
```


### **Event Endpoints**
These endpoints are used for receiving, storing and retrieving events.

### **Create Event Endpoint**
- POST /api/events

This endpoint is used for creating a new event. It accepts a JSON payload containing the event data, and returns the created event data in JSON format.
For the write-intensive application, this endpoint uses rabbitMQ to process DB write operations.

**Request Headers**

Authorization: Bearer [JWT token]

**Request Body**
```
{
  "event_type": <int>,
  "event_time": <string>,
  "user_id": <string>,
  "service_name": <string>,
  "status": <int>,
  "event_fields": {
    "field1": <value>,
    "field2": <value>,
    ...
  }
}
```

**Response:**

   - **Status Code: 200**

   - **Body:**

```
{
  "message": "Successfully added create event to the MQ"
}
```  

### **GetAllEvents Endpoint**
- GET /api/events

This endpoint retrieves all events stored in the system.

**Request Headers**

Authorization: Bearer [JWT token]

**Response:**

   - **Status Code: 200**

   - **Body:**

```
{
    {
        "event_type": <int>,
        "event_time": <string>,
        "user_id": <string>,
        "service_name": <string>,
        "status": <int>,
        "event_fields": {
          "field1": <value>,
          "field2": <value>,
          ...
        }
    },
    ...
}
```
### **GetEventByID Endpoint**
- GET /api/events/:id

This endpoint retrieves event with specific ID stored in the system.

**Request Headers**

Authorization: Bearer [JWT token]

**Response:**

   - **Status Code: 200**

   - **Body:**
```
{
    "event_type": <int>,
    "event_time": <string>,
    "user_id": <string>,
    "service_name": <string>,
    "status": <int>,
    "event_fields": {
      "field1": <value>,
      "field2": <value>,
      ...
    }
}
```

### **GetEventsByTypeID Endpoint**

- GET /api/events/types/:id

This endpoint retrieves events with specific typeIDs stored in the system.

**Request Headers**

Authorization: Bearer [JWT token]

**Response:**

   - **Status Code: 200**

   - **Body:**
```
{
    {
        "event_type": <int>,
        "event_time": <string>,
        "user_id": <string>,
        "service_name": <string>,
        "status": <int>,
        "event_fields": {
          "field1": <value>,
          "field2": <value>,
          ...
        }
    },
    ...
}
```

### **DeleteEventByID Endpoint**
- DELETE /api/events/:id

This endpoint deletes event with specific ID stored in the system.

**Request Headers**

Authorization: Bearer [JWT token]

**Response**

   - **Status Code: 200**

   - **Body:**
```
{
    "message": "Successfully deleted the event"
}
```

## RabbitMQ implementation

This microservice API uses RabbitMQ to handle the creation of new events. RabbitMQ is a message broker that allows applications to communicate with each other asynchronously by sending and receiving messages through a messaging queue. In this implementation, the Create Event endpoint sends a message to a RabbitMQ queue, and a separate consumer service handles the message and stores the new event in the database.

### Architecture
The RabbitMQ implementation consists of the following components:

- **Producer:** The Create Event endpoint is the producer in this implementation. It sends a message to the RabbitMQ queue containing the information about the new event.

- **Queue:** RabbitMQ provides a messaging queue that holds the messages produced by the Create Event endpoint until they are consumed by the consumer service.

- **Consumer:** The consumer service is a separate microservice that is responsible for handling the messages from the RabbitMQ queue. When a new message arrives, the consumer service reads the message and stores the new event in the database.

**Implementation Details**
To implement the RabbitMQ integration, we used the github.com/streadway/amqp library for Go. This library provides a simple API for interacting with RabbitMQ.
