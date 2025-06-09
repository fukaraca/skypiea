FROM golang:1.24-alpine AS builder

ARG FULL_VERSION="dev"
WORKDIR /src
COPY . /src

RUN go build -v -ldflags="-X 'main.Version=${FULL_VERSION}'" -o /src/server ./cmd/server

FROM alpine AS runtime
ARG FULL_VERSION="dev"

LABEL name="skypiea-ai-server"
LABEL org.opencontainers.image.source="https://github.com/fukaraca/skypiea"
LABEL org.opencontainers.image.authors="@fukaraca"
LABEL org.opencontainers.image.description="skypiea server app"
LABEL org.opencontainers.image.version="${FULL_VERSION}"

COPY --from=builder /src/server /app/skypiea-ai/
COPY --from=builder /src/web/ /app/skypiea-ai/web/
COPY --from=builder /src/configs/config.example.yml /app/skypiea-ai/configs/

WORKDIR /app/skypiea-ai
EXPOSE 8080
ENTRYPOINT ["/app/skypiea-ai/server"]
