FROM golang:1.23.4-alpine3.20 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ../common/ ./
COPY assessment.go ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /assessment-serv

FROM alpine:latest
WORKDIR /
COPY ../certs /
COPY --from=build /assessment-serv /assessment-serv
COPY config.json /
EXPOSE 8080
CMD ["/healthguide-serv"]
