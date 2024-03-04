package dac

type DBType string

const (
	Clickhouse DBType = "clickhouse"
	Mysql      DBType = "mysql"
	Postgres   DBType = "postgres"
)
