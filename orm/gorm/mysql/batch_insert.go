package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
	"time"
)

var flowIgnoreFields = []string{"Wall", "Ext", "Loc", "wall", "ext", "loc", "ID", "Id", "DeletedAt"}

// data must be slice
func DoBatchInsert(tableName string, data interface{}, db *gorm.DB) error {
	batch := NewBatchInsertSql(tableName)
	rv := reflect.ValueOf(data)
	rvLen := rv.Len()
	for i := 0; i < rvLen; i++ {
		batch.Add(rv.Index(i).Interface())
	}
	return db.Exec(batch.ResultSql()).Error
}

// 每次批量插入都要new一个
func NewBatchInsertSql(tableName string) *BatchInsertSql {
	nowT := time.Now()
	return &BatchInsertSql{
		TableName: tableName,
		createdAt: nowT.Format("2006-01-02 15:04:05"),
	}
}

// 批量插入组装sql
type BatchInsertSql struct {
	TableName string
	// 保存了属性信息，可以在组装sql时根据属性做不同的操作
	Fields    []reflect.StructField
	InsertSql string

	createdAt string
}

func (b *BatchInsertSql) Add(obj interface{}) {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// 尚未初始化时，初始化sql
	if len(b.Fields) == 0 {
		fieldNames, insertS := getInsertFieldStr(rv.Type(), flowIgnoreFields)
		b.Fields = fieldNames
		b.InsertSql = "insert into " + b.TableName + insertS + " values " + b.getObjValuesForSql(rv, b.Fields)
	} else {
		b.InsertSql += "," + b.getObjValuesForSql(rv, b.Fields)
	}
}

func (b *BatchInsertSql) ResultSql() string {
	return b.InsertSql
}

// 获取插入sql的字段部分
func getInsertFieldStr(rt reflect.Type, ignoreFs []string) (fieldNames []reflect.StructField, fStr string) {
	fStr = "("
	EnumAnObjFieldNames(rt, func(f reflect.StructField) {
		tmpName := f.Name
		// 如果没有ignore则纳入要用的字段中
		if !StrSliceContains(ignoreFs, tmpName) {
			fieldNames = append(fieldNames, f)
			fStr += gorm.ToDBName(tmpName) + ","
			// 保证ID只被加一次
		} else if f.Name == "ID" && f.Type.Kind() == reflect.String {
			fieldNames = append(fieldNames, f)
			fStr += gorm.ToDBName(tmpName) + ","
		}
	})
	if fStr != "" {
		fStr = fStr[:(len(fStr) - 1)]
	}
	fStr += ")"
	return
}

// 获取插入sql的值部分
func (b *BatchInsertSql) getObjValuesForSql(rv reflect.Value, fields []reflect.StructField) (result string) {
	result = "("
	for _, f := range fields {
		//logger.Debug(f.Type.Kind().String())
		//logger.Debug(f.Type.Name())
		// 尚未实现根据类型做适配，因此必须都是string
		if f.Type.Kind() == reflect.Struct && strings.Contains(f.Type.Name(), "Time") {
			if f.Name == "CreatedAt" || f.Name == "UpdatedAt" {
				result += "'" + b.createdAt + "',"
			} else {
				result += "'" + rv.FieldByName(f.Name).Interface().(time.Time).Format("2006-01-02 15:04:05") + "',"
			}
		} else if f.Type.Kind() == reflect.String {
			result += "'" + ClearData4str(rv.FieldByName(f.Name).String()) + "',"
		} else if f.Type.Kind() == reflect.Bool {
			if rv.FieldByName(f.Name).Bool() {
				result += "'1',"
			} else {
				result += "'0',"
			}
		} else {
			result += fmt.Sprintf("'%v',", rv.FieldByName(f.Name).Interface())
		}
	}
	result = result[:(len(result) - 1)]
	result += ")"
	return
}

// 清洗特殊字符串，目前有：
// 1. 单引号转义(批量插入时报错)
func ClearData4str(str string) string {
	if strings.Contains(str, "'") {
		return strings.Replace(str, "'", " ", -1)
	}
	return str
}


// 看一个数组中是否含有某个元素
func StrSliceContains(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

// 迭代一个对象的所有字段名
func EnumAnObjFieldNames(rv reflect.Type, cb func(f reflect.StructField)) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	num := rv.NumField()
	for i := 0; i < num; i++ {
		tmpF := rv.Field(i)
		tmpType := tmpF.Type
		// 如果是时间就不能迭代了
		if tmpType.Kind() == reflect.Struct && !strings.Contains(tmpType.Name(), "Time") && tmpF.Tag.Get("skip") != "true" {
			EnumAnObjFieldNames(tmpType, cb)
		} else {
			cb(tmpF)
		}

	}
}

