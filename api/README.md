# Zopdev API Server

This repository contains the API server for Zopdev. The server is designed to handle backend operations and provide RESTful endpoints for the application.

## Requirements

Ensure you have the following installed on your system before proceeding:

- [Go](https://golang.org/) (latest if possible or v1.20)
- [Docker](https://www.docker.com/) (optional, for containerized development)
- A SQLite database for local development
- [Git](https://git-scm.com/)

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/zopdev/zopdev.git
    cd zopdev/api
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Set up environment variables:
    Create a `.env` file in the root of the project and configure the required environment variables. Example:
    ```env
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_user
    DB_PASSWORD=your_password
    DB_NAME=your_database
    DB_DIALECT=sqlite3
    ```

## Development

1. Start the development server:
    ```bash
    go run main.go
    ```

2. Run tests:
    ```bash
    go test ./...
    ```

3. Use Docker for development (optional):
    Build and run the Docker container:
    ```bash
    docker-compose up --build
    ```

4. API Documentation:
    The API documentation is available at `/api` when the server is running.

## Contributing

Feel free to open issues or submit pull requests to improve the project. Follow the [contribution guidelines](../CONTRIBUTING.md) for more details.
