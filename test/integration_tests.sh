#!/bin/sh

alias grotto=./bin/grotto

# Stop and start PostgreSQL container
docker-compose down
docker-compose up -d

# Give some time for the containers to go up
sleep 2.5

# Build grotto
make

# Test valid migration
grotto -user user -password 123 -database test \
    -dir test/valid_migration

if [[ $? = 0 ]]; then
    echo "Success"
else
    echo "Failure on valid migration test"
    exit -1
fi

# Test migration with syntax error
grotto -user user -password 123 -database test \
    -dir test/migration_with_syntax_error

if [[ $? = 0 ]]; then
    echo "Success"
else
    echo "Failure on migration with syntax error"
    exit -1
fi

# Stop containers
docker-compose down
