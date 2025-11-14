module github.com/kunal768/cmpe202/orchestrator

go 1.23.0

toolchain go1.24.7

require (
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/joho/godotenv v1.5.1
	github.com/kunal768/cmpe202/http-lib v0.0.0
	github.com/rabbitmq/amqp091-go v1.10.0
	go.mongodb.org/mongo-driver v1.14.0
	golang.org/x/crypto v0.41.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
)

replace github.com/kunal768/cmpe202/http-lib => ../http-lib
