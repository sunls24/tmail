FROM oven/bun:alpine AS bun-builder
WORKDIR /app

COPY ./web/package.json ./web/bun.lock ./
RUN bun install --frozen-lockfile
COPY ./web .
RUN bunx astro telemetry disable && bun run build

FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=bun-builder /app/dist ./web/dist/
RUN CGO_ENABLED=0 go build -ldflags '-s -w' -o tmail cmd/main.go

FROM alpine AS runner
WORKDIR /app
COPY --from=builder /app/tmail .
RUN apk add --no-cache tzdata

ENV HOST=127.0.0.1
ENV PORT=3000
EXPOSE 3000
CMD ["/app/tmail"]
