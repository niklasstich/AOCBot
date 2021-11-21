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
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux go build -o aocbot -v .

FROM alpine:latest
COPY --from=build /aocbot/aocbot /app/aocbot
CMD ["/app/aocbot"]
