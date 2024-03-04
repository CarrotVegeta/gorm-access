package dac

import "fmt"

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

func (co *ClickhouseOperator) BuildQuery(condition Condition, qf *QueryFilter) {
	switch condition.Operator {
	case Equal:
		co.Equals(condition, qf)
	case NotEqual:
		co.NotEqual(condition, qf)

	}
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
