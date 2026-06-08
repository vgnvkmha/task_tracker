FROM golang:1.25-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/task-tracker ./cmd/api

FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN useradd --system --uid 10001 --home-dir /app appuser

COPY --from=builder /out/task-tracker /app/task-tracker

RUN chown appuser:appuser /app/task-tracker \
    && chmod 500 /app/task-tracker

USER appuser

EXPOSE 8080

ENV GIN_MODE=release

CMD ["/app/task-tracker"]