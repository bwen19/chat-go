# Go Natter Api

A chat backend api built with Golang

## Getting Started

### Install tools

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [Golang](https://golang.org/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [DBML CLI](https://www.dbml.org/cli/#installation)
- [Sqlc](https://github.com/kyleconroy/sqlc#installation)
- [Gomock](https://github.com/golang/mock)

### Setup infrastructure

- Create the docker network `make network`
- Start postgres container `make postgres`
- Start redis container `make redis`
- Create chat database `make createdb`
- Run db migration up all versions `make migrateup`
- Run db migration down all versions `make migratedown`
- Generate SQL CRUD with sqlc `make sqlc`
- Run Golang server `make server`

## License

Licensed under the [Apache License, Version 2.0](/LICENSE)
