FROM golang:1.23-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /reporadar ./cmd/server

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app
WORKDIR /app

COPY --from=builder /reporadar /usr/local/bin/reporadar
COPY web ./web

USER app
EXPOSE 8080

ENV PORT=8080
ENTRYPOINT ["reporadar"]
