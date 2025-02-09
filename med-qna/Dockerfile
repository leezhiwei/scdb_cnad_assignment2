FROM golang:1.23.4-alpine3.20 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ../common/ ./
COPY signal/ ./
COPY medqna.go ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /medqna-serv

FROM alpine:latest
WORKDIR /
COPY ../certs /
COPY --from=build /medqna-serv /medqna-serv
EXPOSE 8080
CMD ["/medqna-serv"]
