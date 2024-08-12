# Getting Started
The endpoints of the task manager API have been documented along with sample requests and sample responses. For direct access to the API requests, import the collection file located in `/docs/`.

To run the application, go to the root directory of the project and run:
```bash
go run ./Delivery
```
>**Make sure to setup the `.env` file at the root of the project before running the API. Read the details in the next section.**

To check whether the API has started running successfully, make a request to `/ping`.

## Features
### Task API
- Get all tasks
- Get tasks by ID
- Create new tasks
- Update tasks by ID
- Delete tasks by ID

### Auth
- Signup User using username and email
- Login User
- Promote User to Admin

## Project Structure
> Delivery: Contains files related to the delivery layer, handling incoming requests and responses.
- `main.go`: Sets up the HTTP server, initializes dependencies, and defines the routing configuration.
- `controllers/controllers.go`: Handles incoming HTTP requests and invokes the appropriate use case methods.
- `routers/routers.go`: Sets up the routes and initializes the Gin router.


> Domain/: Defines the core business entities and logic.

- `domain.go`: Contains the core business entities such as Task and User structs, the interface definitions for controllers, usecases and repositories along with a customer `error` model used to communicate errors throughout the application

> Infrastructure/: Implements external dependencies and services.
- `auth_middleWare.go`: Middleware to handle authentication and authorization using JWT tokens.
- `jwt_service.go`: Functions to generate and validate JWT tokens.
- `password_service.go`: Functions for hashing and comparing passwords to ensure secure storage of user credentials.

> Repositories/: Abstracts the data access logic.
- task_repository.go: Interface and implementation for task data access operations.

- user_repository.go: Interface and implementation for user data access operations.

> Usecases/: Contains the application-specific business rules.
- task_usecases.go: Implements the use cases related to tasks, such as creating, updating, retrieving, and deleting tasks.
- user_usecases.go: Implements the use cases related to users, such as registering, logging in.

### Tests and Mocks
> Tests/: Contains all the unit tests for the various components of the application

> Mocks/: Contains all the mocked components used in the tests.

## Enviornment Variables

**[IMPORTANT]** There has been a change in how the environment variables are organized. The project now uses the `.env` file located in the root directory with the help of the `viper` package for managing these constants. Additionally, there are additional variables that need to be declared.

The environment variables are as follows:
- `DB_ADDRESS` - connection string of monogoDB
- `SECRET_TOKEN` - used to sign and validate json-web-tokens
- `DB_NAME` - the name of the database instance of the provided connection
- `TEST_DB_NAME` - **[TESTING]** the name of the database instance on which all repository tests will be performed.
- `PORT` - port to run the API on
- `TIMEOUT` - time to wait for operations (in seconds)
- `TOKEN_LIFESPAN_MINUTES` - sets the lifespan of json-web-tokens (in minutes)

**Sample `.env`**
```
DB_ADDRESS=mongodb://localhost:27017
SECRET_TOKEN=long_random_text 
DB_NAME=task_API
TEST_DB_NAME=task_API_test
PORT=8080
TIMEOUT=1
TOKEN_LIFESPAN_MINUTES=30
```
There are no changes to the response and request formats during the transition to using MongoDB for data persistance.

## Sending requests using tokens

The authenitcation system is based on JWT. The token will be sent to the client when it makes a request to `/login` with the correct credentials. That token must be included in the `Authorization` header of any requests to protected routes. The format of the token follows the standard `bearer e...` format. Any deviation from this format will cause the middleware to block the incoming request.

**Sample request (with auth header):**
```bash
curl --location --request DELETE 'http://localhost:8080/tasks/3' \
--header 'Authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Imt5c2sifQ.pIb58jAfa9Rd3u38AzTLdtU_hGR624P6by2epR_baMM'
```

## Running Tests
To run the tests, make sure to first go to the `Test/` directory. **The `repository` tests WILL NOT pass if a valid test database has been setup. Make sure to check the environment variables section for more information on how to provide a test DB.**

To run all the tests:
```bash
go test
```

Each file contains a suite that groups up all the related tests. To run one of suites contained in one of the files, run this command:
```bash
go test -run file_name_test.go NAME_OF_THE_SUITE
```
With timeout:
```bash
go test -timeout 30s -run file_name_test.go NAME_OF_THE_SUITE
```
The suites are usually the functions defined last, accepting a `*testing.T` as a parameter and running the test suite.

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

`http://localhost:8080/login`

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

## Promote

### Authorization: `admin`

**METHOD: PATCH**

`http://localhost:8080/promote/:username`

This endpoint allows admins to promote a user specified by their unique username in the route parameters.

### Response Body

After a successful promote request, the response will be sent with a status code of `200`. The body will contain a JSON object with a message.
- `message`: A message indicating the result of the request to promote a user.

**Example Request (CURL):**
```bash
curl --location --request PATCH 'http://localhost:8080/promote/user1234' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVzQXQiOiIyMDI0LTA4LTEwVDExOjU1OjAzLjY5NTIyODAwNSswMzowMCIsInVzZXJuYW1lIjoiYWRtaW4xMjMifQ.NhVn-7QD67yoT1CQ2ibjzaVTLGuJOIxAqmUerTjDfZ0'
```

**Example Response Body:**
```json
{
    "message": "Used promoted successfully"
}
```

# Task API

## Get Tasks

### Authorization: `user` `admin`

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

### Authorization: `user` `admin`

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

### Authorization: `admin`

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

### Authorization: `admin`

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