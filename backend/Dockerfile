# Multi-stage build: compile Go binary, then run in minimal image
FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# go-fitz uses purego on Linux — no CGO required
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    libffi \
    mupdf \
    mupdf-tools \
    musl-dev

ENV LD_LIBRARY_PATH=/usr/lib

RUN mkdir -p /app/uploads /app/exports

COPY --from=builder /server /app/server

EXPOSE 8080
CMD ["/app/server"]
