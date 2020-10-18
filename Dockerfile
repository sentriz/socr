# syntax=docker/dockerfile:experimental

FROM node:13-buster-slim AS builder-frontend
WORKDIR /socr
COPY frontend/ .
RUN npm install
RUN PRODUCTION=true npm run-script build

FROM golang:1.15-buster AS builder-backend
RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev

WORKDIR /socr
COPY backend/ .
COPY --from=builder-frontend /socr/dist /socr-frontend
RUN	--mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	go get github.com/rakyll/statik
RUN statik -p assets -src /socr-frontend
RUN ls -la /socr-frontend
RUN ls -la
RUN ls -la assets
RUN	--mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	GOOS=linux go build -o socr cmd/socr/main.go

FROM debian:buster-slim
RUN apt-get update -qq
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata/
RUN apt-get install -y -qq tesseract-ocr-eng

COPY --from=builder-backend /socr/socr /
ENV SOCR_LISTEN_ADDR :80
ENV SOCR_SCREENSHOTS_PATH /screenshots
ENV SOCR_INDEX_PATH /index
ENV SOCR_IMPORT_PATH /import
ENTRYPOINT ["/socr"]
