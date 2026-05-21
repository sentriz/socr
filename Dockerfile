# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM alpine:3.23 AS builder-frontend
RUN apk add --no-cache nodejs npm
WORKDIR /src
COPY ./web .
RUN npm install
RUN PRODUCTION=true npm run-script build

FROM --platform=$BUILDPLATFORM golang:1.26-alpine3.23 AS builder-backend
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
COPY --from=builder-frontend /src/dist web/dist/
RUN  \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/ ./cmd/...

FROM alpine:3.23
LABEL org.opencontainers.image.source=https://github.com/sentriz/socr
RUN apk add --no-cache ffmpeg tesseract-ocr tesseract-ocr-data-eng
COPY --from=builder-backend /out/* /usr/local/bin/
ENV SOCR_LISTEN_ADDR=:80
ENV SOCR_DB_DSN=postgres://socr:socr@db:5432?sslmode=disable
ENTRYPOINT [ "/usr/local/bin/socr" ]
