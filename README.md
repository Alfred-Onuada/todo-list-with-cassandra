# Todo List API with Go and Cassandra

This project is a simple RESTful API for managing a to-do list, built with Go and using Apache Cassandra as its database. The API provides basic CRUD (Create, Read, Update, Delete) functionality for to-do items.

## Basic Idea

The application is designed to be a straightforward example of a service that interacts with a NoSQL database. It uses the Gin web framework for routing and handling HTTP requests, and the `gocql` driver to communicate with the Cassandra database.

## Project Structure

- `api/`: Contains the Go source code for the API.
  - `cmd/main.go`: The entry point of the application.
  - `pkg/`: Contains the core logic of the application.
    - `db/`: Handles the database connection and queries.
    - `handlers/`: Contains the HTTP request handlers.
    - `routes/`: Defines the API routes.
  - `Dockerfile`: Used to build the Docker image for the API.
- `cassandra-data/`: Stores the Cassandra database files.
- `docker-compose.yml`: Defines the services, networks, and volumes for the Docker application.

## Setup and Running the Application

To run the application, you need to have Docker and Docker Compose installed on your machine.

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Alfred-Onuada/todo-list-with-cassandra.git
   cd todo-list-with-cassandra
   ```

2. **Start the application:**

   ```bash
   docker compose up --build
   ```

You need to provide and `API_PORT` environment variable
This command will start the API server and the Cassandra database in the foreground. The API will be available at `http://localhost:9990`.

## API Endpoints

The following endpoints are available:

- `GET /`: Health check endpoint.
- `GET /todos`: Get all to-do items.
- `GET /todos/:id`: Get a single to-do item by its ID.
- `POST /todos`: Create a new to-do item.
- `PATCH /todos/:id`: Update an existing to-do item.
- `DELETE /todos/:id`: Delete a to-do item.
