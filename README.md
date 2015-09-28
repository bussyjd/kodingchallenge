# kodingchallenge
This is the answer to a golang challenge

## Requirements
RabbitMQ, Redis, mongodb and postgres must be running:

docker run -d -P --name my_rabbitmq rabbitmq
docker run -d -P --name my_redis redis
docker run -d -P --name my_mongo mongo
docker run -d -P --name my_postgres postgres

## Tests
go test app_test.go app.go

## Build
go build kodingchallenge/worker/{accountname,distinctname,hourlylog}


*Count the occurences of different metrics.
*Return the average occurences of each incoming values
