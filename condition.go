package dac

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

type Operator string

// 操作符常量
const (
	Equal       Operator = "equal"
	NotEqual    Operator = "notEqual"
	GreaterThan Operator = "greaterThan"
	LessThan    Operator = "lessThan"
	Like        Operator = "like"
	In          Operator = "in"
	NotIn       Operator = "notIn"
	IsNull      Operator = "isNull"
	IsNotNull   Operator = "isNotNull"
	NotLike     Operator = "notLike"
	NotBetween  Operator = "notBetween"
	Between     Operator = "between"
)

// OperatorI 定义了操作符接口
type OperatorI interface {
	BuildQuery(condition Condition, query *QueryFilter)
}

var OperatorMap = map[DBType]OperatorI{}

func RegisterOperator(dbType DBType, operator OperatorI) {
	if OperatorMap == nil {
		OperatorMap = make(map[DBType]OperatorI)
	}
	OperatorMap[dbType] = operator
}
func GetOperatorI(dbType DBType) OperatorI {
	return OperatorMap[dbType]
}

// Condition 条件结构体
type Condition struct {
	Field     any // 字段名
	Key       string
	Operator  Operator    // 操作符
	Value     interface{} // 值
	Joiner    JoinerType  // 条件连接词：AND 或 OR
	SubGroups []Condition // 子条件组
}

func NewCondition() *Condition {
	return &Condition{}
}
func (cb *Condition) Build(field any, operator Operator, value any, joiner ...JoinerType) *Condition {
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

func NewConditionBuilder() *ConditionBuilder {
	return &ConditionBuilder{}
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
func (cb *ConditionBuilder) AddCondition(condition Condition) *ConditionBuilder {
	cb.conditions = append(cb.conditions, condition)
	return cb
	//cb.args = append(cb.args, condition.Value)
}

// BuildCondition AddCondition 方法用于添加条件
func (cb *ConditionBuilder) BuildCondition(field any, operator Operator, value any, joiner ...JoinerType) *ConditionBuilder {
	cb.AddCondition(*NewCondition().Build(field, operator, value, joiner...))
	return cb
}

// Append  方法用于添加条件
func (cb *ConditionBuilder) Append(builder ConditionBuilder) {
	cb.conditions = append(cb.conditions, builder.conditions...)
	//cb.args = append(cb.args, condition.Value)
}

// Build 方法用于生成 SQL 条件语句
func (cb *ConditionBuilder) Build(dbType DBType) (string, []interface{}) {
	if cb.conditions == nil {
		return "", nil
	}
	var sqlConditions []string
	for _, condition := range cb.conditions {
		condition.Key = parseField(condition.Field, dbType)
		//value 可能为数组什么的
		condition.Value = parseValue(condition.Value, dbType)
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

// 将条件添加到Where查询中
func addWhereConditions(db *gorm.DB, dbType DBType, builder *ConditionBuilder) error {
	query, args := builder.Build(dbType)
	*db = *db.Where(query, args...)
	return nil
}

// 将条件添加到查询中
func addHavingConditions(db *gorm.DB, dbType DBType, builder *ConditionBuilder) error {
	query, args := builder.Build(dbType)
	*db = *db.Having(query, args...)
	return nil
}
