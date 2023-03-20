# Audit Logging Microservice API

This microservice API is designed to receive, store, and retrieve audit log events. It is built with Golang and MongoDB, and uses RabbitMQ for write-intensive operations.

## Features

- User authentication using email and password
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

### User Endpoints

**User Registration**
- POST /api/users/register

Registers a new user with the provided email and password.

**Request Payload**
```
{
    "email": "user@example.com",
    "password": "mypassword"
}
```
**Response Payload**
```
{
    "id": "6165152c2d9ca44f2c4ee4b1",
}
```

**User Login Endpoint**
- POST /api/users/login

This endpoint is used for authenticating registered users. It accepts a JSON payload containing the user's email and password, and returns a JSON Web Token (JWT) if the email and password are valid.

**Request Body**
```
{
    "email": "user@example.com",
    "password": "password123"
}
```
**Response**

**Status code:** 200 OK

**Body:**

```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE2MjIyNzE5MjYsImlhdCI6MTYyMjI1NTkyNn0.3ey9mCz8WtYr5JrK5rOmLHvzZwWQ2F2ASuoKjOYlGJk"
}
```
Possible Errors

**Status code:** 400 Bad Request

**Body:**

```
{
    "message": "Invalid request payload"
}
```
**Status code:** 401 Unauthorized

**Body:**
```
{
    "message": "Invalid email or password"
}
```
**Status code:** 500 Internal Server Error

**Body:**
```
{
    "message": "Internal server error"
}
```
### Event Endpoints
These endpoints are used for receiving, storing and retrieving events.

**Create Event Endpoint**
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

**Response**

**Status Code: 200**

**Body:**

```
{
  "message": "Successfully added create event to the MQ"
}
```  

**Get All Events Endpoint**
- GET /api/events

This endpoint retrieves all events stored in the system.

**Request Headers**

Authorization: Bearer [JWT token]

**Response**

**Status Code: 200**

**Body:**

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

**GetEventByID**
- GET /api/events/:id

This endpoint retrieves event with specific ID stored in the system.

**Request Headers**

Authorization: Bearer [JWT token]

**Response**

**Status Code: 200**

**Body:**
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

**GetEventsByTypeID**
- GET /api/events/types/:id

This endpoint retrieves events with specific typeIDs stored in the system.

**Request Headers**

Authorization: Bearer [JWT token]

**Response**

**Status Code: 200**

**Body:**
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

**DeleteEventByID**
- DELETE /api/events/:id

This endpoint deletes event with specific ID stored in the system.

**Request Headers**

Authorization: Bearer [JWT token]

**Response**

**Status Code: 200**

**Body:**
```
{
    "message": "Successfully deleted the event"
}
```



