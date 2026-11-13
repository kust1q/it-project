FROM golang:1.27.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./cmd

FROM alpine:3.21

RUN apk add --no-cache ca-certificates && adduser -D -H app
WORKDIR /app

COPY --from=build /app/server ./server
COPY config/config.yml ./config/config.yml

USER app
EXPOSE 8080
ENTRYPOINT ["./server"]
