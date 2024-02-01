package database

import (
	"fmt"
	"gorm.io/gorm"
	"runtime"
)

// DataAccess 数据访问接口
type DataAccess interface {
	Table(db *gorm.DB, name string) DataAccess
	Find(db *gorm.DB, out interface{}, conditions ...interface{}) DataAccess
	First(db *gorm.DB, out interface{}, conditions ...interface{}) DataAccess
	Last(db *gorm.DB, out interface{}, conditions ...interface{}) DataAccess
	Count(db *gorm.DB, count *int64) DataAccess
	Select(db *gorm.DB, fields []Field) DataAccess
	Limit(db *gorm.DB, page, pageSize int64) DataAccess
	Where(db *gorm.DB, conditions ...Condition) DataAccess
}

// 定义全局 map
var registeredDatabases map[DBType]DataAccess

// 在程序启动时注册不同数据库类型的实现
func init() {
	mysqlDB := &MySQLDatabase{}
	RegisterDatabase(Mysql, mysqlDB)

	clickHouseDB := &ClickHouseDatabase{}
	RegisterDatabase(Clickhouse, clickHouseDB)
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

// RegisterDatabase 注册不同数据库类型的方法
func RegisterDatabase(dbType DBType, dataAccess DataAccess) {
	if registeredDatabases == nil {
		registeredDatabases = make(map[DBType]DataAccess)
	}
	registeredDatabases[dbType] = dataAccess
}

// GetDataAccess 方法根据外部传入的数据库类型执行相应的操作
func GetDataAccess(dbType DBType) DataAccess {
	// getDataAccess 函数用于根据 dbType 获取对应的 dataAccess
	if dataAccess, ok := registeredDatabases[dbType]; ok {
		return dataAccess
	}
	return nil
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

func (d *Database) Where(builder ConditionBuilder) *Database {
	err := addConditions(d.db, d.dbType, builder)
	if err != nil {
		d.err = err
	}
	return d
}
func (d *Database) Find(out interface{}) *Database {
	return d.useSourceDB(d.db.Find(out))
}
func (d *Database) Scan(out interface{}) *Database {
	return d.useSourceDB(d.db.Scan(out))
}
func (d *Database) First(out interface{}) *Database {
	return d
}
func (d *Database) Last(out interface{}) *Database {
	return d.useSourceDB(d.db.Last(&out))
}
func (d *Database) Count(count *int64) *Database {
	return d.useSourceDB(d.db.Count(count))
}
func (d *Database) Join(query string, args ...interface{}) *Database {
	d.useSourceDB(d.db.Joins(query, args...))
	return d
}
func (d *Database) Select(fields ...Field) *Database {
	var query string
	for _, v := range fields {
		//验证函数是否合法
		if err := v.validateFunction(); err != nil {
			d.err = err
			return d
		}
		//生成查询 sql
		if query == "" {
			query = v.GenerateSelectSQL()
		} else {
			query += "," + v.GenerateSelectSQL()
		}
	}
	return d.useSourceDB(d.db.Select(query))
}
func (d *Database) Model(model interface{}) *Database {
	return d.useSourceDB(d.db.Model(model))
}
func (d *Database) DB() *gorm.DB {
	return d.db
}
func (d *Database) Limit(page, pageSize int) *Database {
	return d.useSourceDB(d.db.Offset(pageSize * (page - 1)).Limit(pageSize))
}
func (d *Database) Error() error {
	var err error
	if d.err != nil {
		err = d.err
	} else {
		err = d.db.Error
	}
	if err != nil {
		PrintCallerInfo()
	}
	return err
}

// PrintCallerInfo 打印调用者信息
func PrintCallerInfo() {
	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		fmt.Println("Failed to retrieve caller information")
		return
	}

	fmt.Printf("Caller file: %s, line: %d\n", file, line)
}
