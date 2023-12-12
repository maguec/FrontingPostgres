# Fronting Postgres


## Setting up environment

Install go version 1.21 or higher

## running the webserver

```
go run fronting.go 
```

## loading data 

The following command will load 100K records

```
curl -s -X POST http://localhost:8080/load 
```

## running the locust job

```
http://localhost:8099/
```

## enabling caching mode

```
curl -X PATCH http://localhost:8080/config/caching
```

### disabling caching mode (default)

```
curl -X PATCH http://localhost:8080/config/initial
```

