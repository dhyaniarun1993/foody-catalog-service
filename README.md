# foody-catalog-service(In Development)
This project exposes Apis(written in Golang) to manage catalog of the foody system.

To know about the architecture of full product, click [here](https://github.com/dhyaniarun1993/foody-documentation "Foody documentation"). 

## Getting Started

### Manual

#### Prerequisites

1. Golang 
2. Mongodb Server
3. Jaeger(Optional)

#### Clone Repo

Clone the repository at $GOPATH/src/github.com/dhyaniarun1993/.

```sh
$ git clone https://github.com/dhyaniarun1993/foody-catalog-service.git
$ cd foody-catalog-service
```

#### Installing Dependencies

Install dependencies using below command

```sh
$ dep ensure
```

#### Application Configuration

All the configuration that the service needs are passed through Environment variables. For development, environment config can be found at cmd/catalog-server/.env file. After updating the 
variables in env file, export the environment variables using following command

```sh
$ source cmd/catalog-server/.env
```

#### Running the Application

Use the command below to run the application

```sh
$ go run cmd/catalog-server/main.go
```

### Docker

Coming Soon

## API Documentation

Note: This service depends on Nginx and Oauth service to validate the Authorization JWT token and pass the claims(user id, user role, client id) in the header to downstream service. While running in standalone mode, just pass user id, user role and client id in Headers "X-User-Id", "X-User-Role", "X-Client-Id" respectively. 

### APIs

- [x] Restaurant Create, Delete Operations(Only merchants are allowed to perform this operations)
- [x] Get Restaurant Near Me(Only customers are allowed to perform this operations)
- [ ] Get Menu of a Restaurant(Both customer and merchant are allowed to perform this operation)
- [x] Add, Get and Remove Category to restaurant(Only merchants are allowed to perform this operations)
- [x] Add, Get and Delete Product with variant to restaurant and category(Only merchants are allowed to perform this operations)
- [x] Add, Get and Remove variant from restaurant and category(Only merchants are allowed to perform this operations)

Refer to the Api documentation below to know more.

Note: Api schema might change.

Api documentation can be found at [link](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/dhyaniarun1993/foody-catalog-service/master/docs/swagger.yaml "Foody API documentation" )

## Technologies Used

* [Golang](https://golang.org/) - Programming Language to build software
* [MongoDB](https://www.mongodb.com/) - Database to presist information

## Author

Created and maintained by Arun Dhyani
