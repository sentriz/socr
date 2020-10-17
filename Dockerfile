# syntax=docker/dockerfile:experimental

FROM golang:1.15-buster AS builder
RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev

WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
        --mount=type=cache,target=/root/.cache/go-build \
        GOOS=linux go build -o socr cmd/socr/main.go

FROM debian:buster-slim
RUN apt-get update -qq
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata/
RUN apt-get install -y -qq tesseract-ocr-eng

COPY --from=builder /src/socr /

ENV SOCR_LISTEN_ADDR :80
ENV SOCR_SCREENSHOTS_PATH /screenshots
ENV SOCR_INDEX_PATH /index
ENV SOCR_IMPORT_PATH /import

ENTRYPOINT ["/socr"]
