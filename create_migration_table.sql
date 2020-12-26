-- The base table for the migration table.
-- This table is responsible to guarantee that event though the migration
-- may run multiple times, only the not executed scripts will be executed.
create table if not exists golang_migration(
        id bigint generated always as identity primary key,
        script_name varchar constraint uk_script_name unique not null,
        created_at timestamp not null default now()
);
