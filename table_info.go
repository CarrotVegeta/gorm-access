package dac

import "fmt"

// TableInfo 接口定义了获取表名和别名的方法
type TableInfo interface {
	TableName() string
	TableAlias() string
}

// AdditionalTableInfo 接口扩展了TableInfo，添加了获取其他表信息的方法
type AdditionalTableInfo interface {
	TableInfo
	OtherTableInfo() string
}

// TableWithAlias 拼接表名和别名的方法
func TableWithAlias(ti TableInfo, a ...string) string {
	aliasStr := ti.TableAlias()
	if len(a) > 0 && a[0] != "" {
		aliasStr = a[0]
	}
	return fmt.Sprintf("%s AS %s", ti.TableName(), aliasStr)
}
