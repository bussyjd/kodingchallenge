# kodingchallenge

This is the anwer to the challenge: https://gist.github.com/cihangir/c95509f16701a34bd575

## Requirements
RabbitMQ and Redis must be running

## Tests
go test app_test.go app.go

## Bonus
To calculate the average of incoming values I would use MapReduce functions into MongoDB to:
*Count the occurences of different metrics.
*Return the average occurences of each incoming values
