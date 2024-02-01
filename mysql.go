package database

// MySQLDatabase 结构体实现 MySQL 数据库访问方法
type MySQLDatabase struct {
	DataAccess
}
type MysqlOperator struct {
	OperatorI
}

func (m MysqlOperator) BuildQuery(condition Condition, qf *QueryFilter) {
	switch condition.Operator {
	case Equal:
		m.Equal(condition, qf)
	case NotEqual:
		m.NotEqual(condition, qf)
	}
}

func (m MysqlOperator) Equal(condition Condition, qf *QueryFilter) {
	qf.And(condition.Field+" = ?", condition.Value)
}
func (m MysqlOperator) NotEqual(condition Condition, qf *QueryFilter) {
	qf.And(condition.Field+" != ?", condition.Value)
}
