# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    GOBIN=/opt/app/bin \
    go install ./cmd/...

# Run stage
FROM alpine:3.23

RUN apk add --update --no-cache \
    libc6-compat \
    ca-certificates

ENV PATH="/opt/app/bin:${PATH}"

# Keep the `server` binary as the container's default entrypoint, but
# also copy all installed binaries so `client` is available for one-off
# commands or interactive use.
ARG app_bin=server
ENV APP_BIN=$app_bin

COPY --from=builder /opt/app/bin/$app_bin /

ENTRYPOINT ["sh", "-c", "/$APP_BIN"]
