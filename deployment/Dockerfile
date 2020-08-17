FROM ubuntu:18.04

RUN apt-get update

RUN apt-get install -y curl

RUN curl -O https://dl.google.com/go/go1.14.3.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.14.3.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

RUN apt-get install -y npm

RUN apt-get install -y make

RUN npm install -g @vue/cli

RUN mkdir /retext
WORKDIR /retext

COPY pkg ./pkg
COPY cmd ./cmd
COPY Makefile ./Makefile

ENV PATH=$PATH:/usr/local/bin
RUN cd /retext/pkg/www/retext && npm install

COPY ./go.mod .
RUN make backend