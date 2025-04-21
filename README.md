# go-rest

## Description

The project is a simple template for a Golang REST API.

The solution is fully containerized using Docker and Docker compose.

## Deployment

The application can be deployed using Docker or running the app directly.

### Using Docker Compose

Deploy the solution using Docker compose:

```
docker-compose up --build -d
```

### Running the app locally
To run the app directly, you need to have a MongoDB instance running. 

A docker file is provided.

## Configuration
The application can be configured through environment variables. 

More info [here](./docs/deployment/environment.md)

## Project Structure
The project is structured in the following way:
``` 
project-root/
├── cmd/
│   └── app/
│       ├── main.go
│       └── build/
│           └── main.exe
├── docs/                # Project documentation
├── docker-files/        # Docker configuration and related files
├── internal/
│   ├── handlers/
│   │   └── ... # Request handlers
│   ├── middleware/
│   │   └── ... # Middleware logic
│   ├── services/
│   │   └── ... # Business logic/services
│   ├── repositories/
│   │   └── ... # Database operations
│   ├── bootstrap/
│   │   └── ... # Initialization logic
│   └── http/
│       └── endpoints.go # HTTP route definitions
├── pkg/
│   ├── config/          # Application configuration files
│   ├── database/        # DB connection and related logic
│   └── utils/           # Shared utility functions

```

## Notes

- This template is customizable and subject to preference