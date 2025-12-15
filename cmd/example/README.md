# API Example

This example shows how to create a simple API gateway using the `pg` package. In order to
run this example, you will need to have docker installed (since a new PostgreSQL database
will be created for this example):

```bash
go run github.com/mutablelogic/go-pg/example@latest
```

A [dataset](https://www.ssa.gov/oact/babynames/) of about 2.1 million baby names is
ingested using a bulk insert operation. Then a HTTP server is started which allows you to
perform the following operations:

```bash
# Get a list of baby names
curl http://localhost:8080/

# Get a baby name by ID
curl http://localhost:8080/{id}

# Delete a baby name by ID
curl -X DELETE http://localhost:8080/{id}

# Update a baby name by ID
curl -X PATCH http://localhost:8080/{id} -d '{"name": "new name"}'
```
