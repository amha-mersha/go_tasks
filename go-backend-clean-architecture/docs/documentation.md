# Task-Manager API Documentation

Welcome to the Task Manager API documentation. This API is designed using the Clean Architecture principles, which ensures a scalable and maintainable codebase. The API allows users to manage tasks and user roles efficiently.

## Base URL

`http://localhost:<port>/api/v1`

## Endpoints

### 1. User Registration

- **Endpoint:** `/user/register`
- **Method:** `POST`
- **Description:** Allows new users to register. The first user registered will automatically be assigned the `admin` role.
- **Request Body:**
  - **Content-Type:** `application/json`
  - **Body:**
    ```json
    {
      "username": "user123",
      "password": "password123"
    }
    ```
- **Response:**
  - **Status Code:** `201 Created`
  - **Body:**
    ```json
    {
      "message": "User registered successfully"
    }
    ```
  - **Error Response:**
    - **Status Code:** `400 Bad Request`
    - **Body:**
      ```json
      {
        "error": "Invalid input"
      }
      ```
    - **Status Code:** `409 Conflict`
    - **Body:**
      ```json
      {
        "error": "User already exists"
      }
      ```

### 2. User Login

- **Endpoint:** `/user/login`
- **Method:** `POST`
- **Description:** Allows existing users to log in and receive a JWT token for authentication.
- **Request Body:**
  - **Content-Type:** `application/json`
  - **Body:**
    ```json
    {
      "username": "user123",
      "password": "password123"
    }
    ```
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
    ```json
    {
      "token": "jwt_token_here"
    }
    ```
  - **Error Response:**
    - **Status Code:** `401 Unauthorized`
    - **Body:**
      ```json
      {
        "error": "Invalid credentials"
      }
      ```

### 3. Get All Tasks

- **Endpoint:** `/task`
- **Method:** `GET`
- **Description:** Retrieves a list of all tasks. Accessible to both `admin` and `user` roles.
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
    ```json
    [
      {
        "id": "task_id_1",
        "title": "Task 1",
        "description": "Description for Task 1",
        "status": "Pending",
        "priority": "High",
        "due_date": "2024-08-10T00:00:00Z",
        "created_at": "2024-08-01T00:00:00Z",
        "updated_at": "2024-08-01T00:00:00Z"
      }
    ]
    ```
  - **Error Response:**
    - **Status Code:** `401 Unauthorized`
    - **Body:**
      ```json
      {
        "error": "Unauthorized access"
      }
      ```

### 4. Get Task by ID

- **Endpoint:** `/task/:id`
- **Method:** `GET`
- **Description:** Retrieves a single task by its ID. Accessible to both `admin` and `user` roles.
- **Parameters:**
  - **Path Parameter:** `id` (string) - The unique identifier of the task.
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
    ```json
    {
      "id": "task_id_1",
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

### 5. Create a New Task

- **Endpoint:** `/task`
- **Method:** `POST`
- **Description:** Creates a new task. This endpoint is restricted to users with the `admin` role.
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
  - **Error Response:**
    - **Status Code:** `401 Unauthorized`
    - **Body:**
      ```json
      {
        "error": "Unauthorized access"
      }
      ```

### 6. Update an Existing Task

- **Endpoint:** `/task/:id`
- **Method:** `PUT`
- **Description:** Updates an existing task by its ID. This endpoint is restricted to users with the `admin` role.
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

### 7. Delete a Task

- **Endpoint:** `/task/:id`
- **Method:** `DELETE`
- **Description:** Deletes a task by its ID. This endpoint is restricted to users with the `admin` role.
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

### 8. Update User Role

- **Endpoint:** `/user/assign`
- **Method:** `POST`
- **Description:** Allows an admin to update the role of a user. This endpoint only updates the role and does not expose or modify the user's password or other sensitive information.
- **Request Body:**

  - **Content-Type:** `application/json`
  - **Body:**
    ```json
    {
      "username": "user123",
      "role": "admin"
    }
    ```
  - **Fields:**
    - `username` (string): The username of the user whose role is to be updated.
    - `role` (string): The new role to assign to the user (e.g., "admin", "user").

- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**
    ```json
    {
      "username": "user123",
      "role": "admin"
    }
    ```
  - **Error Response:**
    - **Status Code:** `400 Bad Request`
      - **Body:**
        ```json
        {
          "error": "Invalid input"
        }
        ```
    - **Status Code:** `404 Not Found`
      - **Body:**
        ```json
        {
          "error": "User not found"
        }
        ```
    - **Status Code:** `500 Internal Server Error`
      - **Body:**
        ```json
        {
          "error": "Failed to update user role"
        }
        ```

## Authentication

- JWT (JSON Web Token) is used for authentication.
- The `AuthMiddleware` checks the JWT token and verifies the user's role before allowing access to certain routes.

## Roles

- **Admin**: Can create, update, delete tasks, and assign roles to users.
- **User**: Can only view tasks.

## Conclusion

This API provides a robust solution for managing tasks and user roles, with clear separation of concerns and security features implemented through JWT authentication and role-based access control. The use of Clean Architecture ensures the scalability and maintainability of the codebase.
