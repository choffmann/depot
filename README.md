# depot

Inventory service for the warehouse, backed by Postgres.

## Run

    cp .env.example .env
    just up

`just --list` shows all recipes. `just test` needs a running database,
`just db` starts one.
