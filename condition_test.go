package database

import (
	"fmt"
	"testing"
)

func TestConditionBuilder(t *testing.T) {
	conditon := NewCondition()
	conditon.Build("name", Equal, "John", And)
	conditon.SetSubGroups(NewCondition().Build("city", Equal, "New York", Or))
	conditon.SetSubGroups(NewCondition().Build("city", Equal, "Bei jing", Or))
	// 创建 ConditionBuilder 实例
	builder := &ConditionBuilder{}
	// 添加条件
	builder.AddCondition(*conditon)
	// 生成 SQL 条件语句
	sqlConditions, args := builder.Build(Mysql)
	fmt.Println("SQL 条件语句:", sqlConditions)
	fmt.Println(args)
}
