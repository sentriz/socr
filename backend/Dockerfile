# syntax=docker/dockerfile:experimental

FROM golang:1.15-alpine3.12 AS build
RUN apk add --no-cache ca-certificates
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o socr cmd/socr/main.go

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /src/socr /

ENV SOCR_LISTEN_ADDR :80
ENV SOCR_SCREENSHOTS_PATH /screenshots
ENV SOCR_INDEX_PATH /index
ENTRYPOINT ["/socr"]
