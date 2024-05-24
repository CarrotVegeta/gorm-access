package dac

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type FunctionType string

// 常量定义
const (
	MaxFunc           FunctionType = "max"
	MinFunc           FunctionType = "min"
	CountFunc         FunctionType = "count"
	CountNoFieldFunc  FunctionType = "count_no_field"
	CountDistinctFunc FunctionType = "count_distinct"
	AvgFunc           FunctionType = "avg"
	SumFunc           FunctionType = "sum"
	DateFormatFunc    FunctionType = "date_format"
	UpperFunc         FunctionType = "upper"
	LowerFunc         FunctionType = "lower"
	ConcatFunc        FunctionType = "concat"
	LengthFunc        FunctionType = "length" // MySQL中为length，PostgresSQL中为LENGTH
	ToDateTimeFunc    FunctionType = "toDateTime"
	DistinctFunc      FunctionType = "distinct"
)

// 将上述方法常量转换成一个 map
var functionMap = map[FunctionType]string{
	MaxFunc:           "max(%s)",
	MinFunc:           "min(%s)",
	CountFunc:         "count(%s)",
	CountDistinctFunc: "count(distinct %s)",
	AvgFunc:           "avg(%s)",
	SumFunc:           "sum(%s)",
	DateFormatFunc:    "DATE_FORMAT(%s, '%%Y-%%m-%%d %%H:%%i:%%s')",
	UpperFunc:         "upper(%s)",
	LowerFunc:         "lower(%s)",
	ConcatFunc:        "concat(%s)",
	LengthFunc:        "length(%s)",
	ToDateTimeFunc:    "toDateTime(%s)",
	DistinctFunc:      "distinct %s",
}

// GetFunctionSQL 获取 functionMap 中的值
func GetFunctionSQL(function FunctionType) string {
	return functionMap[function]
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

// 定义全局 map
var registeredDataFunctionProvider map[DBType]FunctionProvider

// RegisterFunctionProvider 注册不同数据库类型的函数方法
func RegisterFunctionProvider(dbType DBType, dataAccess FunctionProvider) {
	if registeredDataFunctionProvider == nil {
		registeredDataFunctionProvider = make(map[DBType]FunctionProvider)
	}
	registeredDataFunctionProvider[dbType] = dataAccess
}
func GetDataFunctionProvider(dbType DBType) FunctionProvider {
	if dat, ok := registeredDataFunctionProvider[dbType]; ok {
		return dat
	}
	return &Field{}
}

//通过 functionType和 FunctionProvider 获取对应的 SQL 函数并执行

func GetFunctionHandlerSQL(function FunctionType, fp FunctionProvider) string {
	switch function {
	case MaxFunc:
		return fp.Max()
	case MinFunc:
		return fp.Min()
	case CountFunc:
		return fp.Count()
	case CountDistinctFunc:
		return fp.CountDistinct()
	case AvgFunc:
		return fp.Avg()
	case SumFunc:
		return fp.Sum()
	case DateFormatFunc:
		return fp.DateFormat()
	case UpperFunc:
		return fp.Upper()
	case LowerFunc:
		return fp.Lower()
	case ConcatFunc:
		return fp.Concat()
	case LengthFunc:
		return fp.Length()
	case ToDateTimeFunc:
		return fp.ToDateTime()
	case DistinctFunc:
		return fp.Distinct()
	default:
		return ""
	}
}

// IsFunctionTypeValid 判断是否在该 map 中
func IsFunctionTypeValid(function FunctionType) bool {
	_, ok := functionMap[function]
	return ok
}

// Field 结构表示一个字段，可能包含聚合函数
type Field struct {
	Name     any          // 字段名
	function FunctionType // 聚合函数
	Alias    string       //别名
}

// NewField 创建一个字段，允许不传入聚合函数
func NewField(name string) *Field {
	f := &Field{Name: name}
	return f
}
func (f *Field) As(as string) *Field {
	f.Alias = as
	return f
}

// GenerateSelectSQL 生成 SELECT 语句
func (f *Field) GenerateSelectSQL(fp string) string {
	if f.Name != "" {
		f.Name = convertToSQLFormat(f.Name)
	}
	if f.function != "" {
		f.Name = fmt.Sprintf(fp, f.Name)
	}
	if f.Alias != "" {
		f.Name = fmt.Sprintf("%s as %s", f.Name, f.Alias)
	}
	return fmt.Sprintf("%v", f.Name)
}

// validateFunction 检查聚合函数是否合法
func (f *Field) validateFunction() error {
	if f.function == "" {
		return nil
	}
	if !IsFunctionTypeValid(f.function) {
		return errors.New(string("invalid function: " + f.function))
	}
	return nil
}

func (f *Field) Max() string {
	return GetFunctionSQL(MaxFunc)
}
func (f *Field) Distinct() string {
	return GetFunctionSQL(DistinctFunc)
}
func (f *Field) Min() string {
	return GetFunctionSQL(MinFunc)
}

func (f *Field) Count() string {
	return GetFunctionSQL(CountFunc)
}

func (f *Field) CountDistinct() string {
	return GetFunctionSQL(CountDistinctFunc)
}

func (f *Field) Avg() string {
	return GetFunctionSQL(AvgFunc)
}

func (f *Field) ToDateTime() string {
	return GetFunctionSQL(ToDateTimeFunc)
}

func (f *Field) Sum() string {
	return GetFunctionSQL(SumFunc)
}

func (f *Field) DateFormat() string {
	return GetFunctionSQL(DateFormatFunc)
}

func (f *Field) Upper() string {
	return GetFunctionSQL(UpperFunc)
}

func (f *Field) Lower() string {
	return GetFunctionSQL(LowerFunc)
}

func (f *Field) Concat() string {
	return GetFunctionSQL(ConcatFunc)
}

func (f *Field) Length() string {
	return GetFunctionSQL(LengthFunc)
}

func Max(field string) *Field {
	return &Field{Name: field, function: MaxFunc}
}
func ToDateTime(field string) *Field {
	return &Field{Name: field, function: ToDateTimeFunc}
}
func Min(field string) *Field {
	return &Field{Name: field, function: MinFunc}
}
func Distinct(field string) *Field {
	return &Field{Name: field, function: DistinctFunc}
}
func Count(field any) *Field {
	return &Field{Name: field, function: CountFunc}
}

func CountDistinct(field string) *Field {
	return &Field{Name: field, function: CountDistinctFunc}
}

func Avg(field string) *Field {
	return &Field{Name: field, function: AvgFunc}
}

func Sum(field string) *Field {
	return &Field{Name: field, function: SumFunc}
}

func DateFormat(field string) *Field {
	return &Field{Name: field, function: DateFormatFunc}
}
func Length(field string) *Field {
	return &Field{Name: field, function: LengthFunc}
}

func Upper(field string) *Field {
	return &Field{Name: field, function: UpperFunc}
}
func Lower(field string) *Field {
	return &Field{Name: field, function: LowerFunc}
}

func Concat(fields ...string) *Field {
	concatFields := strings.Join(fields, ", ")
	return &Field{Name: concatFields, function: ConcatFunc}
}

// 解析传进来的参数 field，如果是字符串，直接返回，如果是 *Field 类型，调用 GenerateSelectSQL 方法
// 如果是其他类型，直接返回
func parseField(param interface{}, dbType DBType) string {
	if value, ok := param.(string); ok {
		strs := splitSelectFieldStr(value)
		selectSqlStr := NewSelectStr("")
		for _, v := range strs {
			if !checkFirstLast(v, "`") {
				//检查是否存在 as条件
				// 定义一个正则表达式，匹配"as"前后都是空格的模式
				re := regexp.MustCompile(`\b as \b`)
				// 检查每个字符串是否匹配正则表达式
				// 使用正则表达式找到所有匹配"as"的位置
				matches := re.FindAllStringIndex(v, -1)
				var sqlFormatStr string
				if len(matches) > 0 {
					part1 := v[:matches[0][0]]
					str1 := convertToSQLFormat(part1)
					part2 := v[matches[0][1]:]
					str2 := convertToSQLFormat(part2)
					sqlFormatStr = str1 + " as " + str2
				} else {
					sqlFormatStr = convertToSQLFormat(v)
				}
				selectSqlStr.Join(sqlFormatStr)
				continue
			}
			selectSqlStr.Join(v)
		}
		return selectSqlStr.Value
	}
	if f, ok := param.(*Field); ok {
		provider := GetDataFunctionProvider(dbType)
		return f.GenerateSelectSQL(GetFunctionHandlerSQL(f.function, provider))
	}
	return ""
}
func parseValue(param interface{}, dbType DBType) any {
	if f, ok := param.(*Field); ok {
		provider := GetDataFunctionProvider(dbType)
		return f.GenerateSelectSQL(GetFunctionHandlerSQL(f.function, provider))
	}
	return param
}
