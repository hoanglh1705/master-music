FROM golang:1.20-alpine as builder

RUN apk add --update --no-cache git build-base
RUN mkdir /build
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o bin/svc ./cmd/api

FROM alpine:latest
COPY ./.env ./.env
RUN mkdir swaggerui
COPY ./swaggerui ./swaggerui

COPY --from=builder /build/bin /bin/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8191

ENTRYPOINT ["/bin/svc", "-c", "/config/.env"]
