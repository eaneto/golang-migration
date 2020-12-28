# Grotto

![The blue grotto in Capri](https://upload.wikimedia.org/wikipedia/commons/d/d4/Jakob_Alt_-_Die_Blaue_Grotte_auf_der_Insel_Capri_-_1835-36.jpeg)
> Painting by [Jakob Alt](https://de.wikipedia.org/wiki/Jakob_Alt)

Very basic tool to manage database migrations for PostgreSQL written in go.

## How it works

The program expects a migration directory on the root of the project with sql files.
The files must have the .sql extension and will be ordered before being executed, for
that reason it's a good idea to implement a naming strategy like `V1_XX`, `V2_XX`,
`V3_XX`.

## TODOs

- Unit tests
- Integration tests (with dockerized database)

## Usage

```bash
make
./bin/grotto <user> <password> <database_name> <migration_directory>
```

### With docker compose example

```bash
docker-compose up -d
make
./bin/grotto user 123 test test_migration
```
