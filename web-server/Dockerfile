FROM golang:alpine as build
RUN apk add --no-cache git 

WORKDIR /src

COPY go.sum /src/
COPY go.mod /src/
RUN go mod download 

COPY main.go /src/
COPY controller.go /src/
COPY rabbitmq.go /src/

RUN go build -o app

FROM alpine:3.12

RUN mkdir -p /app
COPY --from=build /src/app /app/app
CMD ["./app/app"]