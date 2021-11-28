# vaccination-demo


## Why?

Learning purpose of trying DDD with Kafka, Rabbit MQ and API in 1 repo
- Kafka publisher using segmentio
- Rabbit MQ use quorum type queue for the purpose of high availability (see consumer_registration.go)
- Standard go-chi for API

## How to run
set .env file
Run `docker compose up -d`

For testing purpose, just curl this example

```
curl --location --request POST 'localhost:8000/resident/1234' \
--header 'Content-Type: application/json' \
--data-raw '{
    "nik":"1234",
    "birth_place":"jakarta",
    "birth_date": "19881212",
    "first_name":"sactio",
    "last_name":"swastioyono"
}'
```

API should response 
```
Sending  {"NIK":"1234","Birthplace":"jakarta","DoB":"19881212","FirstName":"sactio","LastName":"swastioyono"}
```

Consumer should show
```
{"level":"info","ts":"2021-11-28 07:51:58","caller":"cmd/consumer_registration.go:75","msg":" > Received message: %s\n","Body: ":"{\"NIK\":\"1234\",\"Birthplace\":\"jakarta\",\"DoB\":\"19881212\",\"FirstName\":\"sactio\",\"LastName\":\"swastioyono\"}"}
2021/11/28 07:51:58 NIK: 1234
{"level":"info","ts":"2021-11-28 07:52:03","caller":"cmd/consumer_registration.go:87","msg":"acked message"}
```