# Split

Track credits and debts among your friends.
- Split travel expenses
- Split restaurant bills

# Architecture

## Schema Diagram

![schema-diagram](./assets/schema.svg)
*Schema diagram created with [dbdiagram.io](https://www.dbdiagram.io/)*

# Tools

| Tool | Purpose |
|--|--|
| [`migrate`](https://github.com/golang-migrate/migrate) | Schema migration |
| [`sqlc`](https://github.com/sqlc-dev/sqlc) | ORM |

# Notes

### `Makefile`
- Useful `make` commands to set up your dev env

### `migrate`
> *Migrate up* and *migrate down* are included in `Makefile

**Migrate up**
`sslmode=disabled` is necessary as SSL is not enabled on the `pgsql` server by default

```bash
migrate -path db/migrations -database "postgresql://root:password@localhost:5432/split_app?sslmode=disabled" -verbose up
```

### ORMs
`sqlc` was used as an ORM for Go.

