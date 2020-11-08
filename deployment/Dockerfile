FROM ubuntu:18.04

# install git for downloading go dependencies
RUN apt-get update && apt-get install -y git

# major dependencies are vendored in the project and extracted here
COPY third_party /third_party
RUN tar xf /third_party/node-v14.5.0-linux-x64.tar.gz
RUN tar xf /third_party/go1.14.4.linux-amd64.tar.gz

# so we can find npm + go
ENV PATH=$PATH:/node-v14.5.0-linux-x64/bin:/go/bin

# create project dir
RUN mkdir retext

WORKDIR retext

# download the go dependencies before copying in source
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# grab the source
COPY pkg ./pkg
COPY cmd ./cmd

# build the frontend
RUN cd pkg/www/retext && npm install && npm run build

# build the backend
RUN mkdir -p bin && go build -o bin/server cmd/server/main.go

CMD ["./bin/server"]