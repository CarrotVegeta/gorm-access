package database

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type JoinerType string

const (
	And JoinerType = "AND"
	Or  JoinerType = "OR"
)

// Condition 条件结构体
type Condition struct {
	Field     string      // 字段名
	Operator  Operator    // 操作符
	Value     interface{} // 值
	Joiner    JoinerType  // 条件连接词：AND 或 OR
	SubGroups []Condition // 子条件组
}

func NewCondition() *Condition {
	return &Condition{}
}
func (cb *Condition) Build(field string, operator Operator, value string, joiner ...JoinerType) *Condition {
	cb.Field = field
	cb.Operator = operator
	cb.Value = value
	if len(joiner) > 0 {
		cb.Joiner = joiner[0]
	}
	return cb
}
func (cb *Condition) SetSubGroups(subGroups *Condition) *Condition {
	cb.SubGroups = append(cb.SubGroups, *subGroups)
	return cb
}

// ConditionBuilder 用于生成 SQL 条件语句的结构体
type ConditionBuilder struct {
	conditions []Condition
	args       []interface{}
}

// And 方法用于设置 AND 连接词
func (cb *ConditionBuilder) And() *ConditionBuilder {
	return cb.setJoiner("AND")
}

// Or 方法用于设置 OR 连接词
func (cb *ConditionBuilder) Or() *ConditionBuilder {
	return cb.setJoiner("OR")
}

// setJoiner 方法用于设置连接词
func (cb *ConditionBuilder) setJoiner(joiner JoinerType) *ConditionBuilder {
	if len(cb.conditions) > 0 {
		cb.conditions[len(cb.conditions)-1].Joiner = joiner
	}
	return cb
}

// AddCondition 方法用于添加条件
func (cb *ConditionBuilder) AddCondition(condition Condition) {
	cb.conditions = append(cb.conditions, condition)
	//cb.args = append(cb.args, condition.Value)
}

// Append  方法用于添加条件
func (cb *ConditionBuilder) Append(builder ConditionBuilder) {
	cb.conditions = append(cb.conditions, builder.conditions...)
	//cb.args = append(cb.args, condition.Value)
}

// Build 方法用于生成 SQL 条件语句
func (cb *ConditionBuilder) Build(dbType DBType) (string, []interface{}) {
	var sqlConditions []string
	for _, condition := range cb.conditions {
		qf := &QueryFilter{}
		GetOperatorI(dbType).BuildQuery(condition, qf)
		sqlConditions = append(sqlConditions, qf.Query)
		cb.args = append(cb.args, qf.Args...)
		// 处理子条件
		if len(condition.SubGroups) > 0 {
			subConditionBuilder := &ConditionBuilder{}
			for _, subCondition := range condition.SubGroups {
				subConditionBuilder.AddCondition(subCondition)
			}
			subSQLCondition, args := subConditionBuilder.Build(dbType)
			sqlConditions = append(sqlConditions, fmt.Sprintf("(%s)", subSQLCondition))
			cb.args = append(cb.args, args...)
		}
	}
	joinedConditions := strings.Join(sqlConditions, fmt.Sprintf(" %s ", cb.conditions[0].Joiner))
	return fmt.Sprintf("(%s)", joinedConditions), cb.args
}

type Operator string

// 操作符常量
const (
	Equal       Operator = "equal"
	NotEqual    Operator = "notEqual"
	GreaterThan Operator = "greaterThan"
	LessThan    Operator = "lessThan"
	Like        Operator = "like"
)

// OperatorI 定义了操作符接口
type OperatorI interface {
	BuildQuery(condition Condition, query *QueryFilter)
}

var OperatorMap = map[DBType]OperatorI{
	Mysql: MysqlOperator{},
}

func GetOperatorI(dbType DBType) OperatorI {
	return OperatorMap[dbType]
}

// 将条件添加到查询中
func addConditions(db *gorm.DB, dbType DBType, builder ConditionBuilder) error {
	query, args := builder.Build(dbType)
	*db = *db.Where(query, args...)
	return nil
}

// 检查条件是否合法
func checkConditions(o Operator) error {
	switch o {
	case Equal, NotEqual, GreaterThan, LessThan, Like:
	default:
		return fmt.Errorf("不支持的操作符: %s", o)
	}
	return nil
}
