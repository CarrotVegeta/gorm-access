package dac

import "gorm.io/gorm"

// MySQLDatabase 结构体实现 MySQL 数据库访问方法
type MySQLDatabase struct {
	DataAccess
}

// Limit 实现Limit方法
func (m *MySQLDatabase) Limit(db *gorm.DB, page, pageSize int64) *gorm.DB {
	db = db.Limit(int(pageSize)).Offset(int(page * pageSize))
	return db
}

type MysqlOperator struct {
	OperatorI
}

func init() {
	RegisterDatabase(Mysql, &MySQLDatabase{})
	RegisterOperator(Mysql, &MysqlOperator{})
}
func (m MysqlOperator) BuildQuery(condition Condition, qf *QueryFilter) {
	switch condition.Operator {
	case Equal:
		m.Equal(condition, qf)
	case NotEqual:
		m.NotEqual(condition, qf)
	case GreaterThan:
		m.GreaterThan(condition, qf)
	case GreaterThanOrEqual:
		m.GreaterThanOrEqual(condition, qf)
	}
}
func (m MysqlOperator) GreaterThanOrEqual(condition Condition, qf *QueryFilter) {
	qf.And(condition.Key+" >= ?", condition.Value)
}
func (m MysqlOperator) GreaterThan(condition Condition, qf *QueryFilter) {
	qf.And(condition.Key+" > ?", condition.Value)
}
func (m MysqlOperator) LessThanOrEqual(condition Condition, qf *QueryFilter) {
	qf.And(condition.Key+" <= ?", condition.Value)
}

func (m MysqlOperator) Equal(condition Condition, qf *QueryFilter) {
	qf.And(condition.Key+" = ?", condition.Value)
}
func (m MysqlOperator) NotEqual(condition Condition, qf *QueryFilter) {
	qf.And(condition.Key+" != ?", condition.Value)
}
