FROM ubuntu:18.04 as staging

# install git for downloading go dependencies
RUN apt-get update && apt-get install -y git gcc openjdk-8-jre make

# major dependencies are vendored in the project and extracted here
COPY third_party /third_party
RUN tar xf /third_party/node-v14.5.0-linux-x64.tar.gz
RUN tar xf /third_party/go1.14.4.linux-amd64.tar.gz

# so we can find npm + go
ENV PATH=$PATH:/node-v14.5.0-linux-x64/bin:/go/bin

# create project dir
RUN mkdir retext

WORKDIR retext

# copy lock dependency and lock files
COPY go.mod ./
COPY go.sum ./
COPY pkg/www/retext/package.json pkg/www/retext/package.json
COPY pkg/www/retext/package-lock.json pkg/www/retext/package-lock.json

# install dependencies for frontend and backend
RUN go mod download
RUN cd pkg/www/retext && npm install

# grab the source
COPY pkg ./pkg
COPY cmd ./cmd
COPY Makefile ./Makefile
COPY swagger.yaml ./swagger.yaml
# build the frontend

RUN mkdir third_party && mv /third_party/swagger-codegen-cli-3.0.20.jar third_party/swagger-codegen-cli-3.0.20.jar
RUN make client-apis
RUN cd pkg/www/retext && npm run build

# build the backend
RUN mkdir -p bin && go build -o bin/server cmd/server/main.go

FROM ubuntu:18.04 as deploy

COPY --from=staging /retext/bin/server /retext/bin/server
COPY --from=staging /retext/pkg/www/retext/dist /retext/pkg/www/retext/dist

WORKDIR /retext

CMD ["./bin/server"]