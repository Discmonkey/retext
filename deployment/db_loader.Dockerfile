FROM golang:1.14-buster

RUN mkdir retext

WORKDIR retext

COPY cmd/db_setup/main.go ./cmd/db_setup/main.go

COPY pkg/db/postgres/migrations/init/init.sql ./pkg/db/postgres/migrations/init/init.sql
COPY pkg/db/postgres/credentials.go ./pkg/db/postgres/credentials.go
COPY go.mod ./
COPY go.sum ./

RUN go build -o main ./cmd/db_setup/main.go

CMD ["./main", "-init_sql=pkg/db/postgres/migrations/init/init.sql"]



