FROM golang:1.24-alpine AS builder
LABEL name="skypiea-ai-server" \
    maintainer="@fukaraca"
WORKDIR /src
COPY . /src
RUN go build -v -o /src/server ./cmd/server

FROM alpine AS secondbuilder
COPY --from=builder /src/server /app/skypiea-ai/
COPY --from=builder /src/web/ /app/skypiea-ai/web/
COPY --from=builder /src/configs/config.example.yml /app/skypiea-ai/configs/
WORKDIR /app/skypiea-ai
EXPOSE 8080
ENTRYPOINT ["/app/skypiea-ai/server"]
