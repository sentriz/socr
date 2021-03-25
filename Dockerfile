# syntax=docker/dockerfile:experimental

FROM node:13-buster-slim AS builder-frontend
WORKDIR /socr
COPY frontend/ .
RUN npm install
RUN PRODUCTION=true npm run-script build


FROM golang:1.16-buster AS builder-backend
RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev

WORKDIR /socr
COPY . .
COPY --from=builder-frontend /socr/dist ./frontend/dist
RUN	--mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	GOOS=linux go build -o socr backend/cmd/socr/main.go


FROM debian:buster-slim
RUN apt-get update -qq
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata/
RUN apt-get install -y -qq tesseract-ocr-eng wait-for-it

COPY --from=builder-backend /socr/socr /
ENV SOCR_LISTEN_ADDR :80
ENV SOCR_DB_DSN postgres://socr:socr@db:5432?sslmode=disable
ENTRYPOINT [ "/socr" ]
