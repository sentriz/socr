FROM node:16-buster-slim AS builder-frontend
WORKDIR /src
COPY . .
RUN npm install
RUN PRODUCTION=true npm run-script build


FROM golang:1.16-buster AS builder-backend
RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev

WORKDIR /src
COPY . .
COPY --from=builder-frontend /src/dist dist/
RUN GOOS=linux go build -o socr cmd/socr/socr.go


FROM debian:buster-slim
LABEL org.opencontainers.image.source https://github.com/sentriz/socr
RUN apt-get update -qq
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata/
RUN apt-get install -y -qq tesseract-ocr-eng

COPY --from=builder-backend /src/socr /
ENV SOCR_LISTEN_ADDR :80
ENV SOCR_DB_DSN postgres://socr:socr@db:5432?sslmode=disable
ENTRYPOINT [ "/socr" ]
