# Grotto

Very basic tool to manage database migrations for PostgreSQL written in go.

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
