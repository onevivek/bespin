FROM golang:1.16-alpine

WORKDIR /src

RUN apk update && apk upgrade && apk add --no-cache ca-certificates openssh git
RUN update-ca-certificates

# enable access to private Nylas repos
RUN mkdir /root/.ssh/
ADD ssh/id_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"

ENV CGO_ENABLED=0

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
