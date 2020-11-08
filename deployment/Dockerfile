FROM ubuntu:18.04

COPY third_party /third_party
RUN tar xf /third_party/node-v14.5.0-linux-x64.tar.gz
RUN tar xf /third_party/go1.14.4.linux-amd64.tar.gz

# so we can find npm + go
ENV PATH=$PATH:/node-v14.5.0-linux-x64/bin:/go/bin

RUN mkdir retext

WORKDIR retext

COPY go.mod ./
COPY pkg ./pkg

RUN cd pkg/www/retext && npm install && npm run build

RUN mkdir -p bin && go build -o bin/server cmd/server/main.go