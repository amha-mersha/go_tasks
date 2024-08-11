# Task-Manager API Documentation

Welcome to the Task Manager API documentation. This API allows registered users to manage tasks efficiently. The API supports user authentication via JWT and enforces role-based access control to ensure secure and appropriate access to resources.

## Base URL

`http://localhost:<port>/api/v1`

## Authentication and Authorization

### User Roles

- **Admin**: Can create, update, delete, and retrieve tasks.
- **User**: Can only retrieve tasks.
- **Anonymous**: Cannot access any task-related endpoints. Must be registered and logged in.

### JWT Token

- **Authorization Header**: All requests to the protected endpoints must include a valid JWT token in the `Authorization` header.
- **Format**: `Authorization: Bearer <token>`

## Public Endpoints

### 1. Register a New User

- **Endpoint**: `/user/register`
- **Method**: `POST`
- **Description**: Registers a new user. If no users exist in the system, the first user will be an admin.
- **Request Body**:
  - **Content-Type**: `application/json`
  - **Body**:
    ```json
    {
      "username": "newuser",
      "password": "password123",
      "role": "user"
    }
    ```
  - **Notes**: The `role` field can only be "user" during registration. If the database is empty, the first user will automatically become an admin.
- **Response**:
  - **Status Code**: `201 Created`
  - **Body**:
    ```json
    {
      "message": "User registered successfully"
    }
    ```

### 2. Login a User

- **Endpoint**: `/user/login`
- **Method**: `POST`
- **Description**: Logs in a registered user and returns a JWT token.
- **Request Body**:
  - **Content-Type**: `application/json`
  - **Body**:
    ```json
    {
      "username": "existinguser",
      "password": "password123"
    }
    ```
- **Response**:
  - **Status Code**: `200 OK`
  - **Body**:
    ```json
    {
      "token": "<jwt_token>"
    }
    ```

## Protected Endpoints (Require JWT)

### 1. Get All Tasks

- **Endpoint**: `/tasks`
- **Method**: `GET`
- **Description**: Retrieves a list of all tasks. Only accessible by authenticated users (admin or user).
- **Response**:
  - **Status Code**: `200 OK`
  - **Body**:
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
      ...
    ]
    ```
  - **Error Response**:
    - **Status Code**: `401 Unauthorized`
    - **Body**:
      ```json
      {
        "error": "Invalid or missing token"
      }
      ```

### 2. Get Task by ID

- **Endpoint**: `/tasks/:id`
- **Method**: `GET`
- **Description**: Retrieves a single task by its ID. Only accessible by authenticated users (admin or user).
- **Response**:
  - **Status Code**: `200 OK`
  - **Body**:
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
  - **Error Response**:
    - **Status Code**: `404 Not Found`
    - **Body**:
      ```json
      {
        "error": "Task not found"
      }
      ```

### 3. Create a New Task (Admin Only)

- **Endpoint**: `/tasks`
- **Method**: `POST`
- **Description**: Creates a new task. Only accessible by authenticated admins.
- **Request Body**:
  - **Content-Type**: `application/json`
  - **Body**:
    ```json
    {
      "title": "New Task",
      "description": "Description for new task",
      "status": "Pending",
      "priority": "Low",
      "due_date": "2024-08-20T00:00:00Z"
    }
    ```
- **Response**:
  - **Status Code**: `201 Created`
  - **Body**:
    ```json
    {
      "message": "Task created successfully"
    }
    ```
  - **Error Response**:
    - **Status Code**: `401 Unauthorized`
    - **Body**:
      ```json
      {
        "error": "Unauthorized access"
      }
      ```

### 4. Update an Existing Task (Admin Only)

- **Endpoint**: `/tasks/:id`
- **Method**: `PUT`
- **Description**: Updates an existing task by its ID. Only accessible by authenticated admins.
- **Request Body**:
  - **Content-Type**: `application/json`
  - **Body**:
    ```json
    {
      "title": "Updated Task",
      "description": "Updated description for task",
      "status": "In Progress",
      "priority": "Medium",
      "due_date": "2024-08-18T00:00:00Z"
    }
    ```
- **Response**:
  - **Status Code**: `200 OK`
  - **Body**:
    ```json
    {
      "message": "Task updated successfully"
    }
    ```
  - **Error Response**:
    - **Status Code**: `404 Not Found`
    - **Body**:
      ```json
      {
        "error": "Task not found"
      }
      ```
    - **Status Code**: `401 Unauthorized`
    - **Body**:
      ```json
      {
        "error": "Unauthorized access"
      }
      ```

### 5. Delete a Task (Admin Only)

- **Endpoint**: `/tasks/:id`
- **Method**: `DELETE`
- **Description**: Deletes a task by its ID. Only accessible by authenticated admins.
- **Response**:
  - **Status Code**: `200 OK`
  - **Body**:
    ```json
    {
      "message": "Task deleted successfully"
    }
    ```
  - **Error Response**:
    - **Status Code**: `404 Not Found`
    - **Body**:
      ```json
      {
        "error": "Task not found"
      }
      ```
    - **Status Code**: `401 Unauthorized`
    - **Body**:
      ```json
      {
        "error": "Unauthorized access"
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
