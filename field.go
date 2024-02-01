package database

import (
	"errors"
	"fmt"
)

type FunctionType string

// 常量定义
const (
	Max   FunctionType = "max"
	Min   FunctionType = "min"
	Count FunctionType = "count"
	Avg   FunctionType = "avg"
	Sum   FunctionType = "sum"
)

// Field 结构表示一个字段，可能包含聚合函数
type Field struct {
	Name     string       // 字段名
	Function FunctionType // 聚合函数，例如 "max", "min", "count" 等，如果没有则为空字符串
}

// NewField 创建一个字段，允许不传入聚合函数
func NewField(name string, function ...FunctionType) Field {
	f := Field{Name: name}
	if len(function) > 0 {
		f.Function = function[0]
	}
	return f
}

// GenerateSelectSQL 生成 SELECT 语句
func (f Field) GenerateSelectSQL() string {
	if f.Function == "" {
		return f.Name
	}
	return fmt.Sprintf("%s(%s)", f.Function, f.Name)
}

// validateFunction 检查聚合函数是否合法
func (f Field) validateFunction() error {
	switch f.Function {
	case "", Max, Min, Count, Avg, Sum:
		return nil
	default:
		return errors.New(string("invalid function: " + f.Function))
	}
}

type FieldBuilder struct {
	fields []Field
}

func NewFiledBuilder() *FieldBuilder {
	return &FieldBuilder{}
}
func (f *FieldBuilder) BuildField(name string, function ...FunctionType) *FieldBuilder {
	f.fields = append(f.fields, NewField(name, function...))
	return f
}
func (f *FieldBuilder) GetFields() []Field {
	return f.fields
}
