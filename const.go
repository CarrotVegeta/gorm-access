package database

type DBType string

const (
	Clickhouse = "clickhouse"
	Mysql      = "mysql"
	Postgres   = "postgres"
)
