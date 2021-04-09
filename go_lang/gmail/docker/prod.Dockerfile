FROM golang:1.16-alpine AS builder

WORKDIR /src

RUN apk update && apk upgrade && apk add --no-cache ca-certificates git openssh
RUN update-ca-certificates

# enable access to private Nylas repos
ARG DEPLOY_KEY
RUN mkdir /root/.ssh/
RUN echo "$DEPLOY_KEY" > /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN git config --global --add url."git@github.com:".insteadOf "https://github.com/"

ENV CGO_ENABLED=0

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/main .

FROM scratch as bin
COPY --from=builder /bin/main /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/main"]
