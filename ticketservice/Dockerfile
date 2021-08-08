FROM golang:alpine AS builder
WORKDIR /src
COPY . .
RUN go mod init testcontainers-go-e2e/ticketservice
RUN go get
RUN go build -o build/app  app.go

FROM alpine
WORKDIR run
COPY --chown=0:0 --from=builder /src/build /run
RUN apk add --update bash curl && rm -rf /var/cache/apk/*
ENTRYPOINT ["./app"]