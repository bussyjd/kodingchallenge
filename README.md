# kodingchallenge
This is the answer to a golang challenge

## Requirements
RabbitMQ, Redis, mongodb and postgres must be running:

    docker run -d -P --name my_rabbitmq rabbitmq
    docker run -d -P --name my_redis redis
    docker run -d -P --name my_mongo mongo
    docker run -d -P --name my_postgres postgres

## Setup Environment
    cd worker/accountname/ && go get
    cd ../distinctname/ && go get
    cd ../hourlylog/ && go get && cd ../../

## Tests

eg: `go test worker/accountname/accountName_test.go worker/accountname/main.go -postgres_host 192.168.99.100 -postgres_port 32771 -amqp_host 192.168.99.100 -amqp_port 5672 -debug_mode true`

## Build
`go build kodingchallenge/worker/{accountname,distinctname,hourlylog}`

## Run
`go run worker/accountname/main.go -postgres_host 192.168.99.100 -postgres_port 32771 -amqp_host 192.168.99.100 -amqp_port 5672 -debug_mode true`

