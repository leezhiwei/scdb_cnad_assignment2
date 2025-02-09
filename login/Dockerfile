FROM golang:1.23.4-alpine3.20 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY vehicle.go ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /vehicle-serv

FROM alpine:latest
WORKDIR /
COPY --from=build /vehicle-serv /vehicle-serv
COPY start.sh /
EXPOSE 8081
CMD ["/start.sh"]
