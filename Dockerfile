# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM oven/bun:alpine AS bun-builder
WORKDIR /app

COPY ./web/package.json ./web/bun.lock ./
RUN bun install --frozen-lockfile
COPY ./web .
RUN bunx astro telemetry disable && bun run build

FROM --platform=$BUILDPLATFORM golang:1.26.5-alpine AS builder
WORKDIR /app
ARG TARGETOS
ARG TARGETARCH
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=bun-builder /app/dist ./web/dist/
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags '-s -w' -o tmail cmd/main.go

FROM --platform=$BUILDPLATFORM alpine AS alpine-assets
RUN apk add --no-cache ca-certificates tzdata

FROM alpine AS runner
WORKDIR /app
COPY --from=alpine-assets /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine-assets /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/tmail .

ENV HOST=127.0.0.1
ENV PORT=3000
EXPOSE 3000
CMD ["/app/tmail"]
