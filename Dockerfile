FROM golang:latest AS build
WORKDIR /aocbot
COPY . .
RUN go get -d -v .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o aocbot -v .

FROM alpine:latest
COPY --from=build /aocbot/aocbot /app/aocbot
CMD ["/app/aocbot"]
