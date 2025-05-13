# Zopdev API Server

This repository contains the API server for Zopdev. The server is designed to handle backend operations and provide RESTful endpoints for the application.

## Requirements

Ensure you have the following installed on your system before proceeding:

- [Go](https://golang.org/) (v1.24 or later)
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
   DB_NAME=zop.db
   DB_DIALECT=sqlite
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
   docker run -d -p 8000:8000 --name zop-api zopdev/api:v0.2.1

   # To run the UI in a separate container:
   docker run -d -p 3000:3000 -e NEXT_PUBLIC_API_BASE_URL='http://localhost:8000' --name zop-ui zopdev/dashboard:v0.2.1
   ```

## Contributing

Feel free to open issues or submit pull requests to improve the project. Follow the [contribution guidelines](../CONTRIBUTING.md) for more details.
