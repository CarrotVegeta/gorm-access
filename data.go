package dac

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"runtime"
	"strings"
)

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

// 定义全局 map
var registeredDataAccess map[DBType]DataAccess

// RegisterDatabase 注册不同数据库类型的方法
func RegisterDatabase(dbType DBType, dataAccess DataAccess) {
	if registeredDataAccess == nil {
		registeredDataAccess = make(map[DBType]DataAccess)
	}
	registeredDataAccess[dbType] = dataAccess
}

// GetDataAccess 方法根据外部传入的数据库类型执行相应的操作
func GetDataAccess(dbType DBType) DataAccess {
	// getDataAccess 函数用于根据 dbType 获取对应的 dataAccess
	if dat, ok := registeredDataAccess[dbType]; ok {
		return dat
	}
	return nil
}

// Database 结构体定义
type Database struct {
	db     *gorm.DB
	dbType DBType
	da     DataAccess
	err    error
}

// NewDatabase 函数用于创建数据库实例
func NewDatabase(dbType DBType) *Database {
	return &Database{da: GetDataAccess(dbType), dbType: dbType}
}

// Use 传入 db
func (d *Database) Use(db *gorm.DB) *Database {
	d.db = db
	return d
}
func (d *Database) useSourceDB(db *gorm.DB) *Database {
	d.db = db
	return d
}
func (d *Database) Table(name string) *Database {
	return d.useSourceDB(d.db.Table(name))
}

// AutoMigrate 创建表
// AutoMigrate 创建表
func (d *Database) AutoMigrate(dst ...interface{}) error {
	//判断是否为支持的数据类型,如果不支持则返回错误
	for _, v := range dst {
		if err := autoMigrateStruct(d.dbType, reflect.TypeOf(v).Elem()); err != nil {
			return err
		}
	}
	return d.db.AutoMigrate(dst...)
}

// autoMigrateStruct 递归解析结构体
func autoMigrateStruct(dbType DBType, t reflect.Type) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 解析内嵌结构体
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			if err := autoMigrateStruct(dbType, field.Type); err != nil {
				return err
			}
			continue
		}

		tag := field.Tag.Get("gorm")
		if tag == "-" {
			continue
		}

		typeValue := extractTypeFromGormTag(tag)
		if typeValue != "" {
			if !IsDatabaseTypeSupported(typeValue) {
				return fmt.Errorf("Field %s of struct %s has unsupported tag type %s\n", field.Name, t.Name(), typeValue)
			}
			newTag := "gorm:" + ReplaceFieldType(dbType, typeValue)
			field.Tag = reflect.StructTag(newTag)
		} else {
			fieldType := field.Type.Kind().String()
			if field.Type.Kind() == reflect.Slice {
				continue
			}
			if field.Type.Kind() == reflect.Struct {
				fieldType = field.Type.String()
			}
			// 如果没有在 gorm 标签中指定类型，则检查字段类型是否符合常量中的类型
			if !IsConstantTypeSupported(strings.ToLower(fieldType)) {
				return fmt.Errorf("Field %s of struct %s has unsupported type %s\n", field.Name, t.Name(), field.Type)
			}
		}
	}
	return nil
}

// unaliasType 递归展开类型别名
func unaliasType(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Ptr && t.Kind() != reflect.Interface {
		return t
	}
	return unaliasType(t.Elem())
}

// Where 构建查询条件
func (d *Database) Where(buildOption *BuilderOption) *Database {
	d.db = buildWhereConditions(d.db, d.dbType, buildOption)
	return d
}

// Find 查询
func (d *Database) Find(out interface{}) *Database {
	return d.useSourceDB(d.db.Find(out))
}

// Create 创建
func (d *Database) Create(out interface{}) *Database {
	return d.useSourceDB(d.db.Create(out))
}
func (d *Database) Save(out interface{}) *Database {
	return d.useSourceDB(d.db.Save(out))
}

// Updates  根据 `struct` 更新属性，只会更新非零值的字段
func (d *Database) Updates(out interface{}) *Database {
	return d.useSourceDB(d.db.Updates(out))
}

// Update 更新单个列
func (d *Database) Update(column string, value interface{}) *Database {
	return d.useSourceDB(d.db.Update(column, value))
}

// Delete  删除
func (d *Database) Delete(out interface{}) *Database {
	return d.useSourceDB(d.db.Delete(out))
}

// Having having条件查询
func (d *Database) Having(builder *ConditionBuilder) *Database {
	err := addHavingConditions(d.db, d.dbType, builder)
	if err != nil {
		d.err = err
	}
	return d
}

// Scan 将数据输出到指定的结构体
func (d *Database) Scan(out interface{}) *Database {
	return d.useSourceDB(d.db.Scan(out))
}

// First 查询第一条
func (d *Database) First(out interface{}) *Database {
	return d
}

// Last 查询最后一条
func (d *Database) Last(out interface{}) *Database {
	return d.useSourceDB(d.db.Last(out))
}

// Count 查询数量
func (d *Database) Count(count *int64) *Database {
	return d.useSourceDB(d.db.Count(count))
}

// Joins 连接查询
func (d *Database) Joins(query string, args ...interface{}) *Database {
	d.useSourceDB(d.db.Joins(query, args...))
	return d
}
func (d *Database) Join(tableWithAlias, condition string) *Database {
	d.useSourceDB(d.db.Joins("JOIN " + tableWithAlias + " on " + condition))
	return d
}
func (d *Database) Preload(query string, args ...interface{}) *Database {
	d.useSourceDB(d.db.Preload(query, args...))
	return d
}

// Select 查询字段
func (d *Database) Select(fields ...any) *Database {
	var query string
	for _, v := range fields {
		field := parseField(v, d.dbType)
		//生成查询 sql
		if query == "" {
			query = field
		} else {
			query += "," + field
		}
	}
	return d.useSourceDB(d.db.Select(query))
}

// Model 设置模型
func (d *Database) Model(model interface{}) *Database {
	return d.useSourceDB(d.db.Model(model))
}

// DB 获取原始的 db
func (d *Database) DB() *gorm.DB {
	return d.db
}

// Limit 分页
func (d *Database) Limit(page, pageSize int) *Database {
	return d.useSourceDB(d.da.Limit(d.db, int64(page), int64(pageSize)))
}

// Group 分组
func (d *Database) Group(group string) *Database {
	return d.useSourceDB(d.db.Group(group))
}

// Order 排序
func (d *Database) Order(order string) *Database {
	return d.useSourceDB(d.db.Order(order))
}

// Error 获取错误
func (d *Database) Error() error {
	var err error
	if d.err != nil {
		err = d.err
	} else {
		err = d.db.Error
	}
	if err != nil {
		PrintCallerInfo(err)
	}
	return err
}

// PrintCallerInfo 打印调用者信息
func PrintCallerInfo(err error) {
	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		fmt.Println("Failed to retrieve caller information")
		return
	}
	fmt.Printf("Caller file: %s, line: %d", file, line)
}
