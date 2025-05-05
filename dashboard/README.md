# ZopDev

Zop is a comprehensive tool for managing cloud infrastructure. It consists of three main components:

1. **zop-ui**: User interface for managing and monitoring cloud resources.
2. **zop-api**: Backend API service.
3. **zop-cli**: Command-line interface for developers and admins.

---

## Installation

### Prerequisites

- Docker installed on your system.
- Node.js version "^18.18.0 || ^19.8.0 || >= 20.0.0" is required.

---

### Running Locally

#### zop-api

Run the following command to pull and start the Docker image for the zop-api:

```bash
    docker run -d -p 8000:8000 --name zop-api zopdev/api:v0.0.1
```

#### zop-ui

Run the following command to pull and start the Docker image for the zop-ui:

```bash
    docker run -d -p 3000:3000 -e NEXT_PUBLIC_API_BASE_URL='http://localhost:8000' --name zop-ui zopdev/dashboard:v0.0.1
```

> **Note:** The environment variable `NEXT_PUBLIC_API_BASE_URL` is used by zop-ui to connect to the
> zop-api. Ensure that the value matches the API's running base URL.
