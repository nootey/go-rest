D:\Utilities\General\go-migrate\migrate.exe 

migrate create -ext js -dir .\pkg\database\migrations\ create_users

rename to json

## Database
- Have an instance of timescaledb running.
A docker-compose file is provided in `./docker-files/mongo`
Run the docker compose and create an image of the db.

### Migrations

To create a migration, use this command:
```sh
migrate create -ext js -dir .\pkg\database\migrations\ create_users
```

> DISCLAIMER: go-migrate doesn't currently support json migrations, so you'll have to manually rename the `.js` file into `.json

Run the migrations with the following command:
```go 
go run cmd/app/main.go migrate {type}
```

You can also use the provided Makefile:
```shell
make migrate type={type}
```

You can use the following options:
- `up` -> runs all migrations
- `down` -> rolls back migrations
- `fresh` -> drops database and runs all migrations
- `reset` -> resets all migrations

## Seeding

Seeders are handled manually. To create a new seeder, create it in `./databse/seeders/workers`
- Make sure to follow the proper declaration.
```go
func SeederName(ctx context.Context, client *mongo.Client, dbName string) {
	
}
```

Run the seeders with the following command:
```go 
go run cmd/app/main.go seed {type}
```
You can use the following options:
- `basic` -> runs the basic seeders for a fresh rollout
- `full` -> runs all defined seeders, for faking data

<hr> 

Seeders require a .seeder.credentials file in `./pkg/config` to read values to seed.

Currently, these are the required parameters:

```js
ADMIN_EMAIL=
ADMIN_PASSWORD=
```