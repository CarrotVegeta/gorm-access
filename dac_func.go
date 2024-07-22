package dac

import (
	"fmt"
	"strings"
)

//
//var Dac *DacWrapper
//
//type DacWrapper struct {
//	*gorm.DB
//	DBType config.DBType
//}
//
//func Init(db *gorm.DB, dbType config.DBType) {
//	Dac = &DacWrapper{
//		DB:     db,
//		DBType: dbType,
//	}
//}

func GenerateSelectSql(fields ...string) string {
	var newFields []string
	for _, v := range fields {
		// 使用逗号分隔字段
		parts := strings.Split(v, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part) // 去除可能的空白字符
			// 使用 strings.ToLower 来统一处理 AS/as
			lowerPart := strings.ToLower(part)
			aliasIndex := strings.Index(lowerPart, " as ")
			if aliasIndex != -1 {
				// 如果存在别名，只处理别名前的部分
				fieldName := strings.TrimSpace(part[:aliasIndex])
				alias := strings.TrimSpace(part[aliasIndex+4:])
				newFields = append(newFields, QuoteIdentifierIfNeeded(fieldName)+" AS "+alias)
			} else {
				// 没有别名，直接处理整个部分
				newFields = append(newFields, QuoteIdentifierIfNeeded(part))
			}
		}
	}
	return strings.Join(newFields, ",")
}

// QuoteIdentifierIfNeeded 根据数据库类型对标识符进行必要的引号处理
// identifier: 需要处理的标识符
// 返回: 处理后的标识符字符串
func QuoteIdentifierIfNeeded(identifier string) string {
	switch DB.DBType {
	case Mysql:
		// 对于MySQL，不需要做任何处理
		return identifier
	//case config.Dm:
	//	// 对于达梦数据库，用双引号括起标识符
	//	// 注意：这里简单地添加了双引号，但在实际应用中，你可能需要处理标识符中的双引号
	//	// 或者检查是否已经是双引号括起来的
	//	return QuoteAfterDot(identifier)
	default:
		// 如果数据库类型未知，可以返回原始标识符或报错
		return identifier
	}
}

// QuoteAfterDot 如果字符串包含'.'，则在'.'之后的部分添加双引号
// 否则，返回原始字符串
func QuoteAfterDot(s string) string {
	dotIndex := strings.Index(s, ".")
	if dotIndex == -1 {
		// 没有找到'.'，直接返回原始字符串
		return `"` + s + `"`
	}
	// 在'.'之后的部分添加双引号
	return s[:dotIndex+1] + `"` + s[dotIndex+1:] + `"`
}

func GenerateIfElseSql(subQuery string) string {
	switch DB.DBType {
	case Postgres:
		return fmt.Sprintf(
			`CASE
        				WHEN %s THEN true
        				ELSE false
    				END`,
			subQuery,
		)
	case Mysql:
		return fmt.Sprintf(
			`SELECT IF(%s, true, false) `,
			subQuery,
		)
	default:
		return ""
	}
}

func Max(field string, alias ...string) string {
	return buildAliasStr(fmt.Sprintf("max(%v)", field), alias...)
}
func ToDateTime(field string, alias ...string) string {
	return buildAliasStr(fmt.Sprintf("toDateTime(%s)", field), alias...)
}
func Min(field string, alias ...string) string {
	return buildAliasStr(fmt.Sprintf("min(%s)", field), alias...)
}
func Distinct(field string) string {
	return fmt.Sprintf("distinct %s", field)
}
func Count(field any, alias ...string) string {
	return buildAliasStr(fmt.Sprintf("count(%v)", field), alias...)
}
func buildAliasStr(sql string, alias ...string) string {
	if len(alias) > 0 {
		return fmt.Sprintf("%s AS %s", sql, alias[0])
	}
	return sql
}
func Limit(offsetNumber, limitNumber int) string {
	switch DB.DBType {
	case Mysql:
		if offsetNumber == 0 {
			return fmt.Sprintf("LIMIT %d", limitNumber)
		} else {
			return fmt.Sprintf("LIMIT %d,%d", offsetNumber, limitNumber)
		}
	case Postgres:
		if offsetNumber == 0 {
			return fmt.Sprintf("LIMIT %d", limitNumber)
		} else {
			return fmt.Sprintf("LIMIT %d OFFSET %d", limitNumber, offsetNumber)
		}
	}
	return ""
}

func GroupConcat(field string, alias ...string) string {
	switch DB.DBType {
	case Mysql:
		return buildAliasStr(fmt.Sprintf("group_concat(%s Separator ',')", field), alias...)
	case Postgres:
		return buildAliasStr(fmt.Sprintf("string_agg(%s::text,',')", field), alias...)
	}
	return ""
}
func Concat(field string, alias ...string) string {
	return buildAliasStr(fmt.Sprintf("concat(%s)", field), alias...)
}
