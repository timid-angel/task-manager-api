# Getting Started
The endpoints of the task manager API have been documented along with sample requests and sample responses. For direct access to the API requests, import the collection file located in `/docs/`.

To run the application, go to the root directory of the project and run:
```bash
go run .
```

To check whether the API has started running successfully, make a request to `/ping`.

## Enviornment Variables

**[IMPORTANT]** The application uses the connection string defined in `/env.go`. To run the application, you must provide a connection string for the DB. This can either be the address of your local mongod instance or the address of an atlas cluster.

In addition to the database connection string, the API also required a secret key that will be used to sign the json-web-tokens. It is advised to use a randomly generated, long text string. A simply access token could potentially break the authorization system.

**Sample `env.go`**
```go
package main

// this connection string will be made an environment variable upon execution
var DB_URL = "mongodb://localhost:27017"
var JWT_SECRET_TOKEN = "ealksdafin1o*(Uj7sldALKJFnk&s^lf%@98)"
```
There are no changes to the response and request formats during the transition to using MongoDB for data persistance.

## Sending requests using tokens

The authenitcation system is based on JWT. The token will be sent to the client when it makes a request to `/login` with the correct credentials. That token must be included in the `Authorization` header of any requests to protected routes. The format of the token follows the standard `bearer e...` format. Any deviation from this format will cause the middleware to block the incoming request.

**Sample request (with auth header):**
```bash
curl --location --request DELETE 'http://localhost:8080/tasks/3' \
--header 'Authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Imt5c2sifQ.pIb58jAfa9Rd3u38AzTLdtU_hGR624P6by2epR_baMM'
```

# Auth
**Caution:** The API allows the creation of admins without any authorization. This has been done to facilitate proper demonstration. Ideally, the admins would be created before deployment and the route for creating admins would be disabled entirely.

## Signup/Register

### Authorization: None

**METHOD: POST**

`http://localhost:8080/signup`

The `POST /signup` endpoint is used to create a new user account. The request should include the user's email, password, username, and role in the request body.

### Request Body

- `email` (string, required): The email address of the user. Must be a valid email address
    
- `password` (string, required): The password for the user account. Minimum of 8 characters
    
- `username` (string, required): The username chosen by the user. Minimum of 3 characters
    
- `role` (string, required): The role or type of account being created. In this version of the API, there are only two roles: "user" and
"admin". Requests with roles outside of these two will be rejected.    

### Response

Upon succesfuly account creation, a status code of `201` will be sent along with the a message.

**Example Request (CURL):**
```bash
curl --location 'http://localhost:8080/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "natms3@gmail.com",
    "password": "this is a very bad password",
    "username": "kysk1",
    "role": "user"
}'
```

**Example Response Body:**
``` json
{
    "message": "Signup successful"
}
```

## Login

### Authorization: None

**METHOD: POST**

`http://localhost:8080/signup`

This endpoint allows users to authenticate and obtain a token for accessing protected resources.

### Request Body
- `username` (text, required): The username of the user.
    
- `password` (text, required): The password of the user.
        

### Response Body

After a successful login, the response will be sent with a status code of `200`. The body will contain a message and a signed json-web-token that will be used to authorize the user's operations.
- `message`: A message indicating the result of the login attempt.
    
- `token`: A token for accessing protected resources.

**Example Request (CURL):**
```bash
curl --location 'http://localhost:8080/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "kysk",
    "password": "this is a very bad password"
}'
```

**Example Response Body:**
```json
{
    "message": "User logged in successfully",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Imt5c2sifQ.pIb58jAfa9Rd3u38AzTLdtU_hGR624P6by2epR_baMM"
}
```


# Task API

## Get Tasks

### Authorization: None

**METHOD: GET**

`http://localhost:8080/tasks`

This endpoint makes an HTTP GET request to retrieve a list of tasks from the server. The response will be in JSON format and will include an array of task objects, each containing the following properties:
- id (string): The unique identifier for the task.
- title (string): The title or name of the task.
- description (string): A brief description of the task.
- due_date (string): The due date for the task.
- status (string): The current status of the task.

**Example Request (CURL):**
```bash
curl --location 'http://localhost:8080/tasks'
```

**Example Response Body:**
```JSON
[
    {
        "id": "123",
        "title": "Task 1",
        "description": "Complete task 1",
        "due_date": "2022-12-31",
        "status": "in_progress"
    },
    {
        "id": "456",
        "title": "Task 2",
        "description": "Review task 2",
        "due_date": "2022-11-30",
        "status": "pending"
    }
]
```


## Get One Task

### Authorization: None

**METHOD: GET**

`http://localhost:8080/tasks/:id`

This endpoint retrieves the details of a specific task. The structure of the task object is identical to the one described in ***GET Tasks***.

**Example Request (CURL):**
```bash
curl --location 'http://localhost:8080/tasks/4'
```

**Example Response Body:**
```json
{
    "id": "123",
    "title": "Task 1",
    "description": "Complete task 1",
    "due_date": "2022-12-31",
    "status": "in_progress"
}
```

## Create Task

### Authorization: `user` `admin`

**METHOD: POST**

`http://localhost:8080/tasks`

This endpoint allows you to create a new task by sending an HTTP POST request to the specified URL. The request should include a JSON payload in the raw request body, with the following parameters:

### Request Body

| Key | Type | Description |
| --- | --- | --- |
| id | text | The unique identifier for the task. |
| title | text | The title of the task. |
| description | text | A description of the task. |
| due_date | text | The due date for the task. |
| status | text | The status of the task. |

The response to the request will have a status code of 201, indicating that the task has been successfully created. The content type of the response will be in JSON format, and it will include the details of the newly created task, with the same parameters as the request payload.

**Example Request (CURL):**
```bash
curl --location 'http://localhost:8080/tasks' \
--header 'Content-Type: application/json' \
--data '{
    "id": "4",
    "title": "Wash dishes",
    "description": "Just wash the dishes",
    "due_date": "2024-08-05T14:50:56.313532456+03:00",
    "status": "pending"
}'
```

**Example Response Body:**
```json
{
    "id": "4",
    "title": "Wash dishes",
    "description": "Just wash the dishes",
    "due_date": "2024-08-05T14:50:56.313532456+03:00",
    "status": "pending"
}
```

## Update Task

### Authorization: `user` `admin`

**METHOD: PUT**

`http://localhost:8080/tasks/:id`

This endpoint is used to update a specific task identified by its ID. The ID is immutable and won't be updated even if it is present in the request body. The remainder of the fields are, however, mutable and will be updated to the new values if present in the request. The response contains the updated details of the task including its ID, title, description, due date, and status.


**Example Request (CURL):**
```bash
curl --location --request PUT 'http://localhost:8080/tasks/4' \
--header 'Content-Type: application/json' \
--data '{
    "title": "Don'\''t wash the dishes",
    "description": "Go to bed instead"
}'
```

**Example Response Body:**
```json
{
    "id": "4",
    "title": "Don't wash the dishes",
    "description": "Go to bed instead",
    "due_date": "2024-08-05T14:50:56.313532456+03:00",
    "status": "pending"
}
```

## Delete Task

### Authorization: `admin`

**METHOD: DELETE**

`http://localhost:8080/tasks/:id`

This endpoint is used to delete a specific task identified by its ID. It returns a 204: No Content if the task with the provided ID is present and has been deleted successfully.


**Example Request (CURL):**
```bash
curl --location --request DELETE 'http://localhost:8080/tasks/4'
```

**Example Response Body:**
```json
// 204 No Content
```