FROM golang:1.18-alpine3.16 AS build

WORKDIR /build

RUN apk add --no-cache git gcc musl-dev

COPY . .

RUN go build -o ./bin/feed-bot ./cmd/feed-bot

FROM alpine:3.16

WORKDIR /app

COPY --from=build /build/bin/feed-bot /app/

RUN apk add --no-cache ca-certificates && \
    addgroup -S -g 5000 feed-bot && \
    adduser -S -u 5000 -G feed-bot feed-bot && \
    chown -R feed-bot:feed-bot .

USER feed-bot

ENTRYPOINT ["/app/feed-bot"]
