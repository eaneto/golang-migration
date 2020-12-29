# Grotto

![The blue grotto in Capri](https://upload.wikimedia.org/wikipedia/commons/e/eb/Heinrich_Jakob_Fried_-_Die_Blaue_Grotte_auf_Capri.jpg)
> Painting by [Jakob Alt](https://de.wikipedia.org/wiki/Jakob_Alt)

Very basic tool to manage database migrations for PostgreSQL written in go.

## How it works

The program expects a migration directory on the root of the project with sql files.
The files must have the .sql extension and will be ordered before being executed, for
that reason it's a good idea to implement a naming strategy like `V1_XX`, `V2_XX`,
`V3_XX`.


## Usage

```bash
make
./bin/grotto -help
./bin/grotto -user <user> -password <password> -database <database_name> -dir <migration_directory>
```

### With docker compose example

```bash
docker-compose up -d
make
./bin/grotto -user user -password 123 -database test -dir test_migration
```


## Basic integration tests with docker compose

```bash
./test/integration_tests.sql
```
