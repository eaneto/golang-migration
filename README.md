# Grotto

![](https://github.com/eaneto/grotto/workflows/Grotto%20CI/badge.svg)
[![codecov](https://codecov.io/gh/eaneto/grotto/branch/main/graph/badge.svg)](https://codecov.io/gh/eaneto/grotto)

![The blue grotto in Capri](https://upload.wikimedia.org/wikipedia/commons/e/eb/Heinrich_Jakob_Fried_-_Die_Blaue_Grotte_auf_Capri.jpg)
> Painting by [Jakob Alt](https://de.wikipedia.org/wiki/Jakob_Alt)

Basic tool to manage database migrations for PostgreSQL.

## How it works

The program will read all sql files for a given directory and execute
all of them in order. When reading the migration directory, *Grotto*
will order every file by their names, if all the scripts use some kind
of name versioning like,`V1_XX.sql`, `V2_XX.sql`, there won't be any
problems with the execution order, but if they don't match any of
these rules and just have plain text names, like, `create_table_x.sql`
or `create_index_y.sql`, you may run into some trouble if your scripts
must be executed in a different order.

All scripts are executed inside a single transaction, so either all of
the scripts will be executed or none.

## Usage

### Build

```bash
make
```

### Help

```bash
./bin/grotto -help
```

### Run example scripts with docker compose

```bash
docker-compose up -d
make
./bin/grotto -user user -password 123 -database test -dir test/valid_migration
```

## Basic integration tests with docker compose

There is a very simple shell script that runs docker compose, compiles
and runs the grotto test for some scripts under de `test` directory.

```bash
./test/integration_tests.sh
```

## TODOs

- Enhance integration tests with psql validations, like validating a
  table was created, or some data was inserted only once.
