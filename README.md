# Chat service

A chat api build with Go

## Setup local development

### Install tools

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [Golang](https://golang.org/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [DBML CLI](https://www.dbml.org/cli/#installation)
- [Sqlc](https://github.com/kyleconroy/sqlc#installation)
- [Gomock](https://github.com/golang/mock)

### Setup infrastructure

- Create the docker network
- Start postgres container
- Create chat database
- Run db migration up all versions
- Run db migration down all versions
- Generate SQL CRUD with sqlc
- Generate DB mock with gomock
- Create a new db migration
