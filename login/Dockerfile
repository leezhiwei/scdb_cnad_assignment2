FROM golang:1.23.4-alpine3.20 AS build
WORKDIR /app
COPY go.mod go.sum ./
COPY common/ ../common
COPY main.go ./main.go
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /login-serv

FROM alpine:latest
WORKDIR /
COPY certs /certs
COPY --from=build /login-serv /login-serv
COPY config.json /
EXPOSE 8080
CMD ["/login-serv"]
