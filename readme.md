![](.github/socr.png)

## socr

please see example [docker-compose](./docker-compose.yml) and run

```shell
$ # ensure db is up and has the correct database
$ docker-compose up -d socr_db
$ docker-compose exec socr_db psql -U socr -c "create database socr"
 
$ # if database exists, start everything
$ docker-compose up -d
$ docker-compose logs --tail 20 -f main
```
