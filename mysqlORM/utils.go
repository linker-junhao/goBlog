package mysqlORM

import (
	"log"
	"reflect"
	"strings"
	"time"
)

func mysqlDateTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}


func mysqlQueryErrorLog(err error) {
	log.Printf("mysql query error: %s", err.Error())
}

func makeConditionClause(c whereCondition) string {
	conditionClause := " "
	for key, val := range c {
		subCondition := strings.Join(val, " = ? " + key + " ")
		if len(val) == 1{
			subCondition = val[0] + " = ? "
		}
		conditionClause += "(" + subCondition + ") and "
	}
	return conditionClause[:len(conditionClause)-4]
}

// 根据插入字段配置从动态类型中取出值并按顺序构成参数，然后返回，方便之后exec操作做参数
func dynamicModelFieldValues(m interface{}, fields []string) []interface{} {
	var insertValues []interface{}
	for _, field := range fields {
		// 将字段转为帕斯卡写法
		tempFieldName := strings.Replace(field, "_", " ", -1)
		tempFieldName = strings.Title(tempFieldName)
		tempFieldName = strings.Replace(tempFieldName, " ", "", -1)
		insertValues = append(insertValues, reflectValue2BasicType(reflect.ValueOf(m).FieldByName(tempFieldName)))
	}

	return insertValues
}

func reflectValue2BasicType(val reflect.Value) interface{} {
	valK := val.Kind()
	valInterface := val.Interface()
	switch valK {
	case reflect.Int64:
		return valInterface.(int64)
	case reflect.Float64:
		return valInterface.(float64)
	case reflect.String:
		return valInterface.(string)
	case reflect.Bool:
		return valInterface.(bool)
	default:
		log.Print("unsupported type")
		return nil
	}
}

func dynamicModelFieldConditionValues(m interface{}, c whereCondition) []interface{} {
	var conditionValues []interface{}
	for _, subCondition := range c {
		conditionValues = append(conditionValues, dynamicModelFieldValues(m, subCondition)...)
	}
	return conditionValues
}