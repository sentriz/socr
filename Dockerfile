FROM alpine:3.18 AS builder-frontend
RUN apk add --no-cache nodejs npm
WORKDIR /src
COPY ./web .
RUN npm install
RUN PRODUCTION=true npm run-script build

FROM alpine:3.18 AS builder-backend
RUN apk add --no-cache build-base go tesseract-ocr tesseract-ocr-dev leptonica-dev
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
COPY --from=builder-frontend /src/dist web/dist/
RUN GOOS=linux go build -o socr cmd/socr/socr.go

FROM alpine:3.18
LABEL org.opencontainers.image.source https://github.com/sentriz/socr
RUN apk add --no-cache ffmpeg tesseract-ocr-data-eng
COPY --from=builder-backend /src/socr /
ENV SOCR_LISTEN_ADDR :80
ENV SOCR_DB_DSN postgres://socr:socr@db:5432?sslmode=disable
ENTRYPOINT [ "/socr" ]
