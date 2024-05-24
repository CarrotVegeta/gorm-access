# 数据访问组件

## 数据库访问组件设计文档
1. 引言
   数据库访问组件旨在提供一个通用的接口，以屏蔽不同数据库间的差异，使得在不同数据库间切换更加方便，并提供一致的数据访问方式。该组件基于GORM库实现，支持多种主流数据库协议，如MySQL、Oracle、PostgreSQL等。

2. 设计目标
- 提供统一的接口，以封装底层数据库的差异性。
- 支持常见的数据访问操作，如查询、更新、删除等。
- 提供函数方法接口，以支持数据库函数的调用。
- 提供操作符接口，用于构建查询条件。
- 使得在服务中的数据库访问操作更加灵活和可扩展。
- 兼容不同数据库，并保持代码的简洁性和可读性。

3. 数据结构
```
// DataAccess 数据访问接口
type DataAccess interface {
	Table(db *gorm.DB, name string) *gorm.DB
	Find(db *gorm.DB, out interface{}) *gorm.DB
	First(db *gorm.DB, out interface{}) *gorm.DB
	Last(db *gorm.DB, out interface{}) *gorm.DB
	Count(db *gorm.DB, count *int64) *gorm.DB
	Select(db *gorm.DB, fields []Field) *gorm.DB
	Limit(db *gorm.DB, page, pageSize int64) *gorm.DB
	Group(db *gorm.DB, group string) *gorm.DB
	Order(db *gorm.DB, order string) *gorm.DB
}

// FunctionProvider 接口定义了一个获取函数的方法
type FunctionProvider interface {
	Max() string
	Min() string
	Count() string
	CountDistinct() string
	Avg() string
	Sum() string
	DateFormat() string
	Upper() string
	Lower() string
	Concat() string
	Length() string
	ToDateTime() string
	Distinct() string
}

type OperatorI interface {
    BuildQuery(condition Condition, query *QueryFilter)
}
```
4. 接口说明
- DataAccess: 数据访问接口，定义了统一的数据访问操作
- FunctionProvider: 函数方法接口，定义了数据库函数的调用方法
- OperatorI: 操作符接口，用于构建查询条件。
5. 数据库实现   
  实现不同数据库的访问功能，可以根据具体需求分别编写实现。例如，针对MySQL、Oracle、PostgreSQL等数据库，分别实现对应的`DataAccess`接口方法、`FunctionProvider`接口方法和`OperatorI`接口方法。

6. 使用示例
```
   
func TestConditionBuilder(t *testing.T) {
	conditon := NewConditionBuilder()
	conditon.AppendCondition("name", Equal, "John", And)

	sql, qf := conditon.Build(Mysql)
	// 生成 SQL 条件语句
	fmt.Printf("sql :%s", sql)
	fmt.Printf("args:%v", qf)
}

type CountList struct {
	Id uint `json:"id"`
	Count int64 `json:"count"`
}
func TestCreate(t *testing.T) {
	option := NewBuilderOption()
	builder := NewConditionBuilder()
	builder.AppendCondition("app_id", Equal, "APP123")
	builder.AppendCondition("id", Equal, 1)
	option.AppendBuilder(builder)
	var cs []CountList
	err := NewDatabase(Mysql).Use(nil).Where(option).Select("id",Count(1).As("count")).
		Group("id").Find(&cs).Error()
	if err != nil {
		t.Logf(err.Error())
	}
	fmt.Println(cs)
}

   ```
7. 总结
   通过以上设计，我们实现了一个通用的数据库访问组件，可以灵活地支持多种主流数据库，提供了统一的接口以及常见的数据访问操作，同时也提供了函数方法接口以及操作符接口，使得在服务中进行数据库访问更加方便、可扩展和灵活，并且保持了代码的简洁性和可读性。