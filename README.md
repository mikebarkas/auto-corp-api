# Auto Corp

This project is a simple Goalng API that runs in a Docker container alongside a Postgres database.

----

## URL Endpoints
`/json` -- Displays a simple message to verify system is working

`/count` -- Displays how many vehicles are in the database for sale

## Search with Query Params
There are three params available to search vehicles.

`/search?make=Ford` -- Chevrolet, Ford, Subaru, BMW, etc.

`/search?price=10000` -- Show vehicles up to the maximum price

`/search?mileage=50000` -- Show vehicles up to the maximum mileage

## Local Development
Clone repo and docker compose up this should start the Postgres database and seed it with example data.

The API will be available at: `localhost:8080/`

----
### Release Tags

| Tags | Note                                                                 |
|:----:|:---------------------------------------------------------------------|
| 0.4  | Update the README with the URLs                                      |
| 0.3  | Start a Go test workflow action file                                 |
| 0.2  | Adds Postgres container and Docker compose with additional endpoints |
| 0.1  | Single endpoint from the Go app in a container                       |
