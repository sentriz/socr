![](.github/socr.png?v=2)

<https://user-images.githubusercontent.com/6832539/142775925-fc89200c-0bf4-42cc-8099-18fb999288cd.mp4>

## running

please see example [docker-compose](./docker-compose.yml) and run

```shell
$ # ensure db is up and has the correct database
$ docker-compose up -d socr_db
$ docker-compose exec socr_db psql -U socr -c "create database socr"

$ # if database exists, start everything
$ docker-compose up -d
$ docker-compose logs --tail 20 -f main
```

### building from source

requires

- node
- npm
- go (1.19+)

at runtime socr shells out to `ffmpeg` and `tesseract`, so both need to be on `PATH`
(e.g. `apt install ffmpeg tesseract-ocr tesseract-ocr-eng` on debian/ubuntu,
`apk add ffmpeg tesseract-ocr tesseract-ocr-data-eng` on alpine,
`pacman -S ffmpeg tesseract tesseract-data-eng` on arch)

```shell
$ go generate ./web/
$ go install ./cmd/socr/socr.go
$ socr
```
