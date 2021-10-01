FROM node:16-buster-slim AS builder-frontend
WORKDIR /src
COPY ./web .
RUN npm install
RUN PRODUCTION=true npm run-script build


FROM golang:1.17-buster AS builder-backend
RUN apt-get update -qq
RUN apt-get install -y -qq build-essential libtesseract-dev libleptonica-dev libavcodec-dev libavutil-dev libavformat-dev libswscale-dev libgraphicsmagick1-dev

WORKDIR /src
COPY . .
COPY --from=builder-frontend /src/dist web/dist/
RUN GOOS=linux go build -o socr cmd/socr/socr.go


FROM debian:buster-slim
LABEL org.opencontainers.image.source https://github.com/sentriz/socr
RUN apt-get update -qq
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata/

# TODO: static build in previous build step to avoid this
RUN apt-get install -y -qq tesseract-ocr-eng build-essential libtesseract-dev libleptonica-dev libavcodec-dev libavutil-dev libavformat-dev libswscale-dev libgraphicsmagick1-dev

COPY --from=builder-backend /src/socr /
ENV SOCR_LISTEN_ADDR :80
ENV SOCR_DB_DSN postgres://socr:socr@db:5432?sslmode=disable
ENTRYPOINT [ "/socr" ]
