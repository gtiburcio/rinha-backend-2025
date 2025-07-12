# Build Stage
FROM golang:1.24 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./src ./src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rinha-gtiburcio ./src

# Package Stage
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/rinha-gtiburcio .

EXPOSE 3000

ENTRYPOINT ["./rinha-gtiburcio"]
