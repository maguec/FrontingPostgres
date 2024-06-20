# Fronting Postgres

## Introduction

This is a demo application to show how to use Redis and Postgres for lazy lookaside caching.
It generates random customer profiles and stores them in a Postgres database.
Look ups can then be peformed and if the caching config is set the returned data will be cached.

This uses [Locust](https://locust.io/) for load testing and allows the user to see the performance of the application in real time.


## Setting up environment

Install go version 1.21 or higher

## webserver configuration

| Environment Variable | Value     | Use                |
|----------------------|-----------|--------------------|
| PORT                 | 8080      | Web server port    |
| PGHOST               | localhost | Postgres server    |
| PGPORT               | 5432      | Postgres port      |
| PGUSER               | postgres  | Postgres user      |
| PGPASSWORD           | PgDbFTW15 | Postgres password  |
| PGDB                 | profiles  | Postgres database  |
| REDIS_SERVER         | localhost | Redis server       |
| REDIS_PORT           | 6379      | Redis port         |
| DATASIZE             | 100000    | Number of profiles |

## starting the data services

```bash
docker-compose up 
```

## running the webserver

```bash
go run fronting.go 
```

## loading data 

The following command will load 100K records

```bash
curl -s -X POST http://localhost:8080/load 
```

## running the locust job

```
http://localhost:8099/
```

## enabling caching mode

```bash
curl -X PATCH http://localhost:8080/config/caching
```

### disabling caching mode (default)

```bash
curl -X PATCH http://localhost:8080/config/initial
```

