#!/bin/sh

alias grotto=./bin/grotto

start_containers() {
    docker-compose up -d
}

stop_containers() {
    docker-compose down
}

# Stop and start PostgreSQL container
stop_containers
start_containers

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
stop_containers
