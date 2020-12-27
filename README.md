# Grotto

![The blue grotto in Capri](https://lh3.googleusercontent.com/proxy/8TkalBXjwe56lc-OKhMF69rEVhU1HGKPr_ZOvOztRQWO0xw3Ii57QgpbCo3o6yvRRpRqCEj_VYU_viAS_p4mcs7em3gEjZyZuAUGWbY84eFbjvIbkWoMk_fp8Fxx88UYyqMGqjwIDeLSL8fEDwOQRTQAkoQh46mO0Q)
> Painting by [Jakob Alt](https://de.wikipedia.org/wiki/Jakob_Alt)

Very basic tool to manage database migrations for PostgreSQL written in go.

## How it works

The program expects a migration directory on the root of the project with sql files.
The files must have the .sql extension and will be ordered before being executed, for
that reason it's a good idea to implement a naming strategy like V1_XX, V2_XX, V3_XX...

## TODOs

- Unit tests
- Integration tests (with dockerized database)

## Usage

```bash
make
./bin/grotto <user> <password> <database_name>
```

### With docker compose example

```bash
docker-compose up -d
make
./bin/grotto user 123 test
```
