# go-rest

## Description

The project is a simple template for a Golang REST API.

The solution is fully containerized using Docker and Docker compose.

## Deployment

The application can be deployed using Docker or running the app directly.

It also supports auto reload with Air.

### Using Docker Compose

Deploy the solution using Docker compose:

```
docker-compose up --build -d
```

### Running the app locally
To run the app directly, you need to have a MongoDB instance running. 

You can also use 'air' and run it with hot reload.

## Configuration
The application can be configured through environment variables. 
The following options can be configured:

```
MONGO_URI=mongodb://root:root@localhost:27017/go-rest?authSource=admin
MONGO_DB=go-rest 
PORT=8080
RELEASE=local
```

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
- API should be secured
- Validation should be added