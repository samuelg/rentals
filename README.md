# Rentals API

## Functionality

This project is a rentals JSON API that returns a list of rentals that can be filtered,
sorted, and paginated implemented in Golang.

The data for the application is stored in a POSTGIS database, running inside docker
by default.

The application supports the following endpoints.

- `/rentals/<RENTAL_ID>` Read one rental endpoint
- `/rentals` Read many (list) rentals endpoint
  - Supported query parameters
    - price_min (number)
    - price_max (number)
    - limit (number)
    - offset (number)
    - ids (comma separated list of rental ids)
    - near (comma separated pair [lat,lng])
    - sort (string)
  - Examples:
    - `rentals?price_min=9000&price_max=75000`
    - `rentals?limit=3&offset=6`
    - `rentals?ids=3,4,5`
    - `rentals?near=33.64,-117.93` // within 100 miles
    - `rentals?sort=price`
    - `rentals?near=33.64,-117.93&price_min=9000&price_max=75000&limit=3&offset=6&sort=price`

The rental object JSON in the response has the following structure:

```json
{
  "id": "int",
  "name": "string",
  "description": "string",
  "type": "string",
  "make": "string",
  "model": "string",
  "year": "int",
  "length": "decimal",
  "sleeps": "int",
  "primary_image_url": "string",
  "price": {
    "day": "int"
  },
  "location": {
    "city": "string",
    "state": "string",
    "zip": "string",
    "country": "string",
    "lat": "decimal",
    "lng": "decimal"
  },
  "user": {
    "id": "int",
    "first_name": "string",
    "last_name": "string"
  }
}
```

When listing rentals the response has the following structure:

```json
{
  "pagination": {
    "count": 30,
    "limit": 10,
    "offset": 0
  },
  "data": [{
    "id": "int",
    "name": "string",
    "description": "string",
    "type": "string",
    "make": "string",
    "model": "string",
    "year": "int",
    "length": "decimal",
    "sleeps": "int",
    "primary_image_url": "string",
    "price": {
      "day": "int"
    },
    "location": {
      "city": "string",
      "state": "string",
      "zip": "string",
      "country": "string",
      "lat": "decimal",
      "lng": "decimal"
    },
    "user": {
      "id": "int",
      "first_name": "string",
      "last_name": "string"
    }
  }]
}
```

## Development

### Requirements

- Golang 1.20
- Docker (to run application inside Docker container)
- Docker compose (to initialize and run the POSTGIS database)

### Running the application

To run the application locally outside of a container:

```sh
go mod download
docker-compose up
go run .
```

To run the application locally inside of a docker container:

```sh
docker-compose -f docker-compose-with-app.yml build
docker-compose -f docker-compose-with-app.yml up
```

### Running tests

To run tests:

```sh
docker-compose up
go test -v ./...
```

## Prometheus metrics

Prometheus metrics for gin routes can be found [here](http://localhost:8080/metrics).

