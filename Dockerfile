# syntax=docker/dockerfile:1.3
FROM golang:latest AS build
ENV CGO_ENABLED=0
#dependencies
WORKDIR /aocbot
COPY go.mod .
COPY go.sum .
RUN go mod download

#build
COPY . .
RUN GOOS=linux go build -o aocbot -v .

FROM ubuntu:latest
RUN apt update
RUN apt install -y fonts-firacode inkscape ca-certificates
ENV SSL_CERT_DIR=/etc/ssl/certs
COPY --from=build /aocbot/aocbot /app/aocbot
CMD ["/app/aocbot"]
