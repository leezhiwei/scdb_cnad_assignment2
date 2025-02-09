FROM golang:1.23.4-alpine3.20 AS build
WORKDIR /app
COPY go.mod go.sum ./
COPY common/ ../common
COPY signal/ ./signal
COPY main.go ./main.go
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /video-serv

FROM alpine:latest
WORKDIR /
COPY certs /certs
COPY --from=build /video-serv /video-serv
EXPOSE 8080
CMD ["/video-serv"]
