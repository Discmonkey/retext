FROM golang:1.14-buster as staging

RUN mkdir retext

WORKDIR retext

# copy go mod and get dependencies before building everything
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/db_setup/main.go ./cmd/db_setup/main.go
COPY pkg/store/postgres_backend/connection.go ./pkg/store/postgres_backend/connection.go
COPY pkg/store/credentials/credentials.go ./pkg/store/credentials/credentials.go
COPY pkg/version/version.go ./pkg/version/version.go

RUN go build -o /main ./cmd/db_setup/main.go

FROM ubuntu:18.04 as deploy

COPY --from=staging /main /main

CMD ["/main", "-migration_dir=/pkg/store/postgres_backend/migrations"]
