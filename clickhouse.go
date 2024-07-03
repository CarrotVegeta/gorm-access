package dac

import (
	"fmt"
	"gorm.io/gorm"
)

func init() {
	RegisterDatabase(Clickhouse, &ClickHouseDatabase{})
}

// ClickHouseDatabase 结构体实现 ClickHouse 数据库访问方法
type ClickHouseDatabase struct {
	DataAccess
}
type ClickhouseOperator struct {
	OperatorI
}

// Limit 实现Limit方法
func (m *ClickHouseDatabase) Limit(db *gorm.DB, page, pageSize int64) *gorm.DB {
	db = db.Limit(int(pageSize)).Offset(int(page * pageSize))
	return db
}

func (co *ClickhouseOperator) BuildQuery(condition Condition, qf *QueryFilter) {
	switch condition.Operator {
	case Equal:
		co.Equals(condition, qf)
	case NotEqual:
		co.NotEqual(condition, qf)
	case GreaterThanOrEqual:
		co.GreaterThanOrEqual(condition, qf)
	case GreaterThan:
		co.GreaterThan(condition, qf)
	case LessThanOrEqual:
		co.LessThanOrEqual(condition, qf)
	case LessThan:
		co.LessThan(condition, qf)
	case In:
		co.In(condition, qf)
	}
}
func (co *ClickhouseOperator) In(condition Condition, qf *QueryFilter) {
	qf.And(condition.Key+" in (?)", condition.Value)
}
func (co *ClickhouseOperator) GreaterThanOrEqual(condition Condition, qf *QueryFilter) {
	qf.And(condition.Key+" >= ?", condition.Value)
}
func (co *ClickhouseOperator) GreaterThan(condition Condition, qf *QueryFilter) {
	qf.And(condition.Key+" > ?", condition.Value)
}
func (co *ClickhouseOperator) LessThanOrEqual(condition Condition, qf *QueryFilter) {
	qf.And(condition.Key+" <= ?", condition.Value)
}
func (co *ClickhouseOperator) LessThan(condition Condition, qf *QueryFilter) {
	qf.And(condition.Key+" < ?", condition.Value)
}
func (co *ClickhouseOperator) NotBlank(condition Condition, qf *QueryFilter) {
	qf.And(fmt.Sprintf("%s <> ''  ", condition.Field), condition.Value)
}
func (co *ClickhouseOperator) Blank(condition Condition, qf *QueryFilter) {
	qf.And(fmt.Sprintf("%s = ''  ", condition.Field), condition.Value)
}

func (co *ClickhouseOperator) NotEqual(condition Condition, qf *QueryFilter) {
	qf.And(fmt.Sprintf("%s <> ? ", condition.Field), condition.Value)
}

// Equals 等于
func (co *ClickhouseOperator) Equals(condition Condition, qf *QueryFilter) {
	qf.And(fmt.Sprintf("%s = ? ", condition.Field), condition.Value)
}
