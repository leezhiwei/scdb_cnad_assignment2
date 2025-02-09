FROM golang:1.23.4-alpine3.20 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /static-serv
FROM alpine:latest
WORKDIR /
COPY static/ /static
COPY --from=build /static-serv /static-serv
EXPOSE 80
CMD ["/static-serv"]
