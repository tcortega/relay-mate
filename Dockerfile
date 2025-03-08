FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN apk --no-cache add sqlite sqlite-dev build-base gcc
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o relay-mate cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/relay-mate .
RUN apk --no-cache add ca-certificates sqlite

ENTRYPOINT ["./relay-mate"]
