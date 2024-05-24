package dac

import (
	"fmt"
	"testing"
)

func TestConditionBuilder(t *testing.T) {
	conditon := NewConditionBuilder()
	conditon.AppendCondition("name", Equal, "John", And)

	sql, qf := conditon.Build(Mysql)
	// 生成 SQL 条件语句
	fmt.Printf("sql :%s", sql)
	fmt.Printf("args:%v", qf)
}

type CountList struct {
	Id    uint  `json:"id"`
	Count int64 `json:"count"`
}

func TestCreate(t *testing.T) {
	option := NewBuilderOption()
	builder := NewConditionBuilder()
	builder.AppendCondition("app_id", Equal, "APP123")
	builder.AppendCondition("id", Equal, 1)
	option.AppendBuilder(builder)
	var cs []CountList
	err := NewDatabase(Mysql).Use(nil).Where(option).Select("id", Count(1).As("count")).
		Group("id").Find(&cs).Error()
	if err != nil {
		t.Logf(err.Error())
	}
	fmt.Println(cs)
}
