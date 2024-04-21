package dac

import (
	"regexp"
	"strings"
)

const (
	TYPE_INT       = "int"    // 整数类型
	TYPE_BIGINT    = "bigint" // 长整型
	TYPE_TINYINT   = "tinyint"
	TYPE_DECIMAL   = "decimal"   // 十进制数
	TYPE_NUMERIC   = "numeric"   // 数值
	TYPE_REAL      = "real"      // 实数
	TYPE_DOUBLE    = "double"    // 双精度浮点数
	TYPE_SMALLINT  = "smallint"  // 短整型
	TYPE_BOOLEAN   = "boolean"   // 布尔型
	TYPE_CHAR      = "char"      // 字符串（定长）
	TYPE_VARCHAR   = "varchar"   // 字符串（变长）
	TYPE_DATE      = "date"      // 日期
	TYPE_TIME      = "time"      // 时间
	TYPE_TIMESTAMP = "timestamp" // 时间戳
	TYPE_INTERVAL  = "interval"  // 时间间隔
	TYPE_BYTEA     = "bytea"     // 二进制数据
	TYPE_UUID      = "uuid"      // UUID
	//不受其他数据库支持的字段类型
	TYPE_TEXT      = "text" // 文本
	TYPE_LONG_TEXT = "longtext"
	TYPE_ENUM      = "enum"
	TYPE_BLOB      = "blob"
)

func ReplaceFieldType(dbType DBType, fieldType string) string {
	switch dbType {
	case Mysql:
		switch fieldType {
		case "text":
			return "customtype"
		}
	case Clickhouse:
		return ""
	}
	return ""
}

// IsDatabaseTypeSupported 检查数据库类型是否受支持
func IsDatabaseTypeSupported(fieldType string) bool {
	databaseTypes := []string{
		TYPE_INT, TYPE_BIGINT, TYPE_DECIMAL, TYPE_NUMERIC, TYPE_REAL,
		TYPE_DOUBLE, TYPE_SMALLINT, TYPE_BOOLEAN, TYPE_CHAR, TYPE_VARCHAR,
		TYPE_TEXT, TYPE_DATE, TYPE_TIME, TYPE_TIMESTAMP, TYPE_INTERVAL,
		TYPE_BYTEA, TYPE_UUID, TYPE_TINYINT, TYPE_LONG_TEXT, TYPE_ENUM, TYPE_BLOB,
	}

	// 匹配 dbType 中的类型名称部分，忽略括号和冒号后面的内容
	re := regexp.MustCompile(`^(\w+)[(:]`)
	match := re.FindStringSubmatch(fieldType)
	if len(match) >= 2 {
		fieldType = match[1]
	}
	for _, databaseType := range databaseTypes {
		if strings.EqualFold(fieldType, databaseType) {
			return true
		}
	}

	return false
}

const (
	GO_TYPE_UINT          = "uint"          // 无符号整数类型
	GO_TYPE_UINT8         = "uint8"         // 无符号 8 位整数类型
	GO_TYPE_UINT16        = "uint16"        // 无符号 16 位整数类型
	GO_TYPE_UINT32        = "uint32"        // 无符号 32 位整数类型
	GO_TYPE_UINT64        = "uint64"        // 无符号 64 位整数类型
	GO_TYPE_INT           = "int"           // 整数类型       -> int
	GO_TYPE_INT8          = "int8"          // 8 位整数类型    -> int8
	GO_TYPE_INT16         = "int16"         // 16 位整数类型   -> int16
	GO_TYPE_INT32         = "int32"         // 32 位整数类型   -> int32
	GO_TYPE_INT64         = "int64"         // 64 位整数类型   -> int64
	GO_TYPE_FLOAT32       = "float32"       // 单精度浮点数   -> float32
	GO_TYPE_FLOAT64       = "float64"       // 双精度浮点数   -> float64
	GO_TYPE_BOOLEAN       = "bool"          // 布尔型         -> bool
	GO_TYPE_STRING        = "string"        // 文本           -> string
	GO_TYPE_TIME          = "time.Time"     // 时间           -> time.Time 或者自定义的 Time 类型
	GO_TYPE_INTERVAL      = "time.Duration" // 时间间隔       -> time.Duration 或者自定义的 Interval 类型
	GO_TYPE_BYTEA         = "[]byte"        // 二进制数据     -> []byte
	GO_TYPE_UUID          = "string"        // UUID           -> string 或者自定义的 UUID 类型
	GO_TYPE_SQL_NULL_TIME = "sql.NullTime"  // SQL 空时间 -> sql.NullTime
)

// extractTypeFromGormTag 从 gorm 标签中提取 type 值
func extractTypeFromGormTag(tag string) string {
	parts := strings.Split(tag, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "type:") {
			return strings.Trim(strings.TrimPrefix(part, "type:"), " ")
		}
	}
	return ""
}

// IsConstantTypeSupported 检查常量类型是否受支持
func IsConstantTypeSupported(constantType string) bool {
	constants := []string{
		GO_TYPE_INT, GO_TYPE_UINT, GO_TYPE_UINT8, GO_TYPE_UINT16, GO_TYPE_UINT32,
		GO_TYPE_UINT64,
		GO_TYPE_INT,
		GO_TYPE_INT8,
		GO_TYPE_INT16,
		GO_TYPE_INT32,
		GO_TYPE_INT64,
		GO_TYPE_FLOAT32,
		GO_TYPE_FLOAT64,
		GO_TYPE_BOOLEAN,
		GO_TYPE_STRING, GO_TYPE_TIME, GO_TYPE_INTERVAL, GO_TYPE_BYTEA,
		GO_TYPE_UUID,
		GO_TYPE_SQL_NULL_TIME,
	}

	for _, constant := range constants {
		if strings.EqualFold(constantType, constant) {
			return true
		}
	}

	return false
}
