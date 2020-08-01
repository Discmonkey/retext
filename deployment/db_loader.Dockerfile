FROM golang:1.14-buster

RUN mkdir retext

WORKDIR retext

# copy go mod and get dependencies before building everything
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/db_setup/main.go ./cmd/db_setup/main.go
COPY pkg/store/credentials/credentials.go ./pkg/store/credentials/credentials.go

RUN go build -o main ./cmd/db_setup/main.go

CMD ["./main", "-init_sql=pkg/db/migrations/init/init.sql"]
