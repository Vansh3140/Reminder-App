# Reminder App

## Overview

The Reminder App is a robust and efficient backend service for managing user authentication and event reminders. Built using Go, Fiber, and MySQL, this project emphasizes secure user management and CRUD operations for events, leveraging JWT-based authentication.

---

## Features

- **User Authentication**: Signup and login functionalities with password hashing using `bcrypt`.
- **JWT Middleware**: Secure API endpoints with JWT token authentication.
- **Event Management**: Create, retrieve, update, and delete events tied to user accounts.
- **Database Integration**: Uses MySQL with proper schema management and foreign key relationships.
- **TLS Support**: Secure database connections using TLS.

---

## Project Structure

```
Reminder-App/
├── main.go          # Application entry point
├── handlers/
│   └── handlers.go  # Event-related logic and API handlers
├── database/
│   └── database.go  # Database connection and schema setup
├── .env             # Environment variables (e.g., DB credentials, JWT secret)
```

---

## Setup Instructions

### Prerequisites

- Go 1.18+
- MySQL Server
- Git

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/Vansh3140/Reminder-App.git
   cd Reminder-App
   ```

2. Set up environment variables:
   Create a `.env` file in the root directory and add the following:
   ```env
   DB_CREDS="username:password@tcp(127.0.0.1:3306)/reminderapp"
   SECRET_KEY="your_secret_key"
   CERTIFICATE="your_tls_certificate"
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

The application will start on `http://localhost:8080`.

---

## API Endpoints

### **Public Endpoints**

#### 1. `POST /signup`
   **Description**: Create a new user account.

   **Request Body**:
   ```json
   {
       "username": "example_user",
       "password": "example_password"
   }
   ```

   **Response**:
   ```json
   {
       "token": "<JWT_TOKEN>"
   }
   ```

#### 2. `POST /login`
   **Description**: Log in and retrieve a JWT token.

   **Request Body**:
   ```json
   {
       "username": "example_user",
       "password": "example_password"
   }
   ```

   **Response**:
   ```json
   {
       "token": "<JWT_TOKEN>"
   }
   ```

### **Protected Endpoints** (Require JWT Token)

Add the JWT token to the `Authorization` header as: `Bearer <JWT_TOKEN>`.

#### 3. `POST /api/v1/event`
   **Description**: Create a new event.

   **Request Body**:
   ```json
   {
       "name": "Meeting",
       "date": "2025-01-15",
       "message": "Team sync-up meeting"
   }
   ```

   **Response**:
   ```json
   {
       "status": "created",
       "event_name": "Meeting",
       "message": "Event created successfully"
   }
   ```

#### 4. `GET /api/v1/event/:name`
   **Description**: Retrieve event details by name.

   **Response**:
   ```json
   {
       "status": "fetched",
       "event_id": 1,
       "details": {
           "name": "Meeting",
           "date": "2025-01-15",
           "message": "Team sync-up meeting"
       },
       "message": "Event fetched successfully"
   }
   ```

#### 5. `PUT /api/v1/event/:name`
   **Description**: Update an event's details.

   **Request Body**:
   ```json
   {
       "name": "Updated Meeting",
       "date": "2025-01-16",
       "message": "Updated team sync-up meeting"
   }
   ```

   **Response**:
   ```json
   {
       "status": "updated",
       "event_id": 1,
       "message": "Event updated successfully"
   }
   ```

#### 6. `DELETE /api/v1/event/:name`
   **Description**: Delete an event by name.

   **Response**:
   ```json
   {
       "status": "deleted",
       "event_name": "Meeting",
       "message": "Event deleted successfully"
   }
   ```

---

## Database Schema

### Users Table
```sql
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);
```

### Events Table
```sql
CREATE TABLE IF NOT EXISTS events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    date VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (name, user_id)
);
```

---

## Security Features

1. **Password Hashing**: User passwords are hashed using `bcrypt` before storing in the database.
2. **JWT Authentication**: Secure token-based authentication for protected routes.
3. **TLS Connection**: Ensures secure database communication with a custom TLS configuration.

---

## Contribution Guidelines

1. Fork the repository.
2. Create a feature branch: `git checkout -b feature-name`.
3. Commit your changes: `git commit -m "Add new feature"`.
4. Push to the branch: `git push origin feature-name`.
5. Open a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

