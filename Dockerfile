FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN apk --no-cache add sqlite sqlite-dev --virtual build-dependencies build-base gcc
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o relay-mate cmd/main.go

FROM debian:stable-slim

WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/relay-mate .

CMD ["./relay-mate"]