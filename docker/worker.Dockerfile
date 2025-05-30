FROM golang:1.24-alpine AS builder
LABEL name="skypiea-ai-worker" \
    maintainer="@fukaraca"
WORKDIR /src
COPY . /src
RUN go build -v -o /src/worker ./cmd/worker

FROM alpine AS secondbuilder
COPY --from=builder /src/worker /app/skypiea-ai/
COPY --from=builder /src/configs/config.example.yml /app/skypiea-ai/configs/
WORKDIR /app/skypiea-ai
ENTRYPOINT ["/app/skypiea-ai/worker"]