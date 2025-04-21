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
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   └── ... # Handlers
│   │   └── middleware/
│   │       └── api.go
│   ├── server/
│   │   ├── server.go
│   │   └── endpoints.go
│   ├── models/
│   ├── repositories/
│   │   ├── mongo/
│   │   │   └── main.go
│   │   └── ... # Other repositories
│   ├── utils/
│   │   └── ... # Helper functions
└── pkg/
    ├── config/         # Application configuration files
    └── ...             # Compiled proto files
```

## Notes

- This template is customizable and subject to preference