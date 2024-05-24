package dac

import (
	"fmt"
	"strings"
)

func convertToSQLFormat(input any) string {
	if _, ok := input.(string); !ok {
		return fmt.Sprintf("%d", input)
	}
	parts := strings.Split(input.(string), ".")
	if len(parts) != 2 {
		return fmt.Sprintf("`%s`", input)
	}
	return fmt.Sprintf("`%s`.`%s`", parts[0], parts[1])
}
func checkFirstLast(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	return s[:len(substr)] == substr && s[len(s)-len(substr):] == substr
}

func splitSelectFieldStr(fieldStr string) []string {
	return strings.Split(fieldStr, ",")
}

// SelectStr 是一个自定义类型，用于字符串的拼接操作。
type SelectStr struct {
	Value string
}

// NewSelectStr 构造函数，用于创建一个新的 SelectStr 实例。
func NewSelectStr(str string) *SelectStr {
	return &SelectStr{Value: str}
}

// Join 方法用于将当前实例与另一个字符串进行拼接。
// 如果当前实例为空，则直接使用另一个字符串；如果不为空，则将其值逗号分隔拼接到当前实例之后。
func (s *SelectStr) Join(str string) *SelectStr {
	if s.Value == "" {
		s.Value = str
		// 如果当前实例为空，直接使用另一个字符串
		return s
	}
	// 如果当前实例不为空，将其与另一个字符串的值以逗号分隔拼接
	s.Value += "," + str
	return s
}
