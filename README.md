# Project Title

RESTful API with Go Chi, GORM, and Custom API Gateway

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Project Structure](#project-structure)
- [Services](#services)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Overview
![Todo Diagram](https://ibb.co/xzcFTLC)
This project is a RESTful API built using the Go Chi router and GORM ORM, encapsulated in Docker containers managed by Docker Compose. It includes an API gateway developed from scratch that handles user authentication and routes requests to the appropriate services (User and Todo services). Each service follows a structured layout to maintain organization and separation of concerns.

## Features

- **Go Chi**: A lightweight, idiomatic, and composable router for building Go HTTP services.
- **GORM**: The fantastic ORM library for Golang.
- **Docker & Docker Compose**: Containerization and orchestration of services.
- **Custom API Gateway**: Built from scratch to manage user authentication and routing.
- **Authentication Service**: Ensures secure access to the API.
- **User and Todo Services**: Core services with CRUD functionality.
- **Comprehensive Testing**: Includes unit tests for all major components.

## Project Structure

The project is organized into multiple services, each following a similar layout:
```
service
├── configs/
├── cmd/app/main.go
├── internal/
│ ├── handlers/
│ ├── routing/
│ ├── tests/
│ ├── errors/
│ ├── middlewares/
│ ├── db/
│ ├── server/
```


## Services

### API Gateway
- **Purpose**: Handle authentication and route requests to the appropriate service.
### Auth Service
- **Purpose**: Manage user authentication.

### User Service
- **Purpose**: Manage user data and operations.

### Todo Service
- **Purpose**: Manage todo items and operations.
- 
## Getting Started

### Prerequisites

- Docker

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/ennemli/todo.git
   cd todo
2. Build and Run:
   ```sh
   docker compose build
   docker compose up
3. Access the API Gateway at http://localhost:8000.

# Usage
## API Endpoints
### User Service:

- POST api/users: Create a new user.
- GET api/users/{id}: Get a user by ID.
- PUT api/users/{id}: Update a user by ID.
- DELETE api/users/{id}: Delete a user by ID.
### Todo Service:

- POST /todos: Create a new todo.
- GET /todos/{id}: Get a todo by ID.
- PUT /todos/{id}: Update a todo by ID.
- DELETE /todos/{id}: Delete a todo by ID.
### Auth Service:

- POST /auth/login: User login.
- POST /auth/valid: User validation.

## Testing
Each service includes unit and integration tests located in the internal/tests/ directory. To run the tests, use the following command:

sh
```sh
docker exec <service-name> go test ./internal/tests/...
```
## License
This project is licensed under the MIT License - see the [MIT License](https://opensource.org/licenses/MIT).

