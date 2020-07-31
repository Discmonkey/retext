FROM golang:1.14-buster

RUN mkdir retext

WORKDIR retext

COPY cmd/db_setup/main.go ./cmd/db_setup/main.go

COPY pkg/db/migrations/init/init.sql ./pkg/db/migrations/init/init.sql
COPY pkg/db/credentials/credentials.go ./pkg/db/credentials/credentials.go

COPY go.mod ./
COPY go.sum ./

RUN go build -o main ./cmd/db_setup/main.go

CMD ["./main", "-init_sql=pkg/db/migrations/init/init.sql"]
