# setting up and building the app
FROM golang:1.24 AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./netzero ./cmd/main.go

# Setting up container that will run the app
FROM alpine:latest AS run

WORKDIR /app

COPY --from=build /build/netzero ./netzero

# COPY .env ./ #TODO: figure out a better solution here

COPY sql ./sql

EXPOSE 8080

#TODO: for tls traffic

# EXPOSE 8081 

CMD ["./netzero"]

