
build:
  go build -tags='no_postgres no_mysql no_vertica no_mssql no_clickhouse no_ydb no_libsql' -o goose ./cmd/goose

