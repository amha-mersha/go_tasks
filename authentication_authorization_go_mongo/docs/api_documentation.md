# Task-Manager

Welcome to the Task Manager API documentation. This API allows users to manage tasks efficiently by providing endpoints to create, read, update, and delete tasks. The API is designed to handle task-related operations such as setting task priorities, tracking due dates, and updating task statuses. The API is built using Go with the Gin web framework and MongoDB as the database. The following endpoints are available to interact with tasks.

## Base URL

`http://localhost:<port>/api/v1`

## Endpoints

### 1. Get All Tasks

- **Endpoint:** `/tasks`
- **Method:** `GET`
- **Description:** Retrieves a list of all tasks.
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
    ```json
    [
      {
        "title": "Task 1",
        "description": "Description for Task 1",
        "status": "Pending",
        "priority": "High",
        "due_date": "2024-08-10T00:00:00Z",
        "created_at": "2024-08-01T00:00:00Z",
        "updated_at": "2024-08-01T00:00:00Z"
      },
      {
        "title": "Task 2",
        "description": "Description for Task 2",
        "status": "Completed",
        "priority": "Medium",
        "due_date": "2024-08-15T00:00:00Z",
        "created_at": "2024-08-02T00:00:00Z",
        "updated_at": "2024-08-05T00:00:00Z"
      }
    ]
    ```

### 2. Get Task by ID

- **Endpoint:** `/tasks/:id`
- **Method:** `GET`
- **Description:** Retrieves a single task by its ID.
- **Parameters:**
  - **Path Parameter:** `id` (string) - The unique identifier of the task.
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
    ```json
    {
      "title": "Task 1",
      "description": "Description for Task 1",
      "status": "Pending",
      "priority": "High",
      "due_date": "2024-08-10T00:00:00Z",
      "created_at": "2024-08-01T00:00:00Z",
      "updated_at": "2024-08-01T00:00:00Z"
    }
    ```
  - **Error Response:**
    - **Status Code:** `404 Not Found`
    - **Body:**
      ```json
      {
        "error": "Task not found"
      }
      ```

### 3. Create a New Task

- **Endpoint:** `/tasks`
- **Method:** `POST`
- **Description:** Creates a new task.
- **Request Body:**
  - **Content-Type:** `application/json`
  - **Body:**
    ```json
    {
      "title": "New Task",
      "description": "Description for new task",
      "status": "Pending",
      "priority": "Low",
      "due_date": "2024-08-20T00:00:00Z"
    }
    ```
- **Response:**
  - **Status Code:** `201 Created`
  - **Body:**
    ```json
    {
      "message": "Task created successfully"
    }
    ```

### 4. Update an Existing Task

- **Endpoint:** `/tasks/:id`
- **Method:** `PUT`
- **Description:** Updates an existing task by its ID.
- **Parameters:**
  - **Path Parameter:** `id` (string) - The unique identifier of the task.
- **Request Body:**
  - **Content-Type:** `application/json`
  - **Body:**
    ```json
    {
      "title": "Updated Task",
      "description": "Updated description for task",
      "status": "In Progress",
      "priority": "Medium",
      "due_date": "2024-08-18T00:00:00Z"
    }
    ```
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
    ```json
    {
      "message": "Task updated successfully"
    }
    ```
  - **Error Response:**
    - **Status Code:** `404 Not Found`
    - **Body:**
      ```json
      {
        "error": "Task not found"
      }
      ```

### 5. Delete a Task

- **Endpoint:** `/tasks/:id`
- **Method:** `DELETE`
- **Description:** Deletes a task by its ID.
- **Parameters:**
  - **Path Parameter:** `id` (string) - The unique identifier of the task.
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
    ```json
    {
      "message": "Task deleted successfully"
    }
    ```
  - **Error Response:**
    - **Status Code:** `404 Not Found`
    - **Body:**
      ```json
      {
        "error": "Task not found"
      }
      ```

## Models

### Task Model

```go
type Task struct {
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    Priority    string    `json:"priority"`
    DueDate     time.Time `json:"due_date"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Title: The title of the task.

Description: A brief description of the task.

Status: The current status of the task (e.g., Pending, In Progress, Completed).

Priority: The priority level of the task (e.g., Low, Medium, High).

DueDate: The due date for the task.

CreatedAt: The timestamp when the task was created.

UpdatedAt: The timestamp when the task was last updated.

## Database

- MongoDB is used to store the tasks.
- Connection is handled in the data/db_connection.go file.

## Environment Variables

- GODOTENV is used for managing environment variables.
- Typical .env file includes:

```bash
MONGO_URI=mongodb+srv://<username>:<password>@cluster0.mongodb.net/taskmanager?retryWrites=true&w=majority
PORT=8080
```

## Running the Application

- Clone the repository.
- Create a .env file in the root directory with the necessary environment variables.
- Run the application:

```bash
go run main.go
```

- The server will be up and running on http://localhost:<port>.

## Conclusion

This API allows users to manage tasks by providing endpoints for creating, reading, updating, and deleting tasks. The app is built using Go, Gin, and MongoDB, providing a scalable solution for task management.

### Notes:

- **Endpoints** are documented with their respective HTTP methods and expected request/response formats.
- **Models** section details the structure of the Task model.
- **Database** and **Environment Variables** sections provide necessary configuration details.

This documentation gives a clear overview of how to use the API and what to expect from each endpoint.
