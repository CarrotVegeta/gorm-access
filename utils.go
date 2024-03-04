package dac

import (
	"fmt"
	"strings"
)

func convertToSQLFormat(input string) string {
	parts := strings.Split(input, ".")
	if len(parts) != 2 {
		return input // 如果格式不正确，直接返回原始字符串
	}
	return fmt.Sprintf("`%s`.`%s`", parts[0], parts[1])
}
func checkFirstLast(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	return s[:len(substr)] == substr && s[len(s)-len(substr):] == substr
}
