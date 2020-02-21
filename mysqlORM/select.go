package mysqlORM

import (
	"log"
	"reflect"
	"strings"
)

type Select struct {
	opBasic
	selectFields    []string
	selectWhereFields whereCondition
	whereModel interface{}
}


func (s *Select)SetSelectFields(f []string) *Select {
	s.setStmtNil()
	s.selectFields = f
	return s
}

func (s *Select)SetSelectWhereFields(calculateType string, fields []string) *Select {
	s.setStmtNil()
	s.selectWhereFields[calculateType] = fields
	return s
}

func (s *Select)Where(model interface{}) *Select {
	s.whereModel = model
	return s
}

func (s *Select)prepareStmt(){
	if s.stmt == nil {
		selectFStr := "*"
		if len(s.selectFields) != 0 {
			selectFStr = strings.Join(s.selectFields, ",")
		}
		theSql := "select " + selectFStr + " from " + s.table + " where " + makeConditionClause(s.selectWhereFields)
		stmt, err := s.db.Prepare(theSql)
		if err != nil {
			log.Print("Update prepare failed:"+err.Error())
		}
		s.stmt = stmt
	}
}

func (s *Select)Commit() ([]interface{}, error) {
	s.prepareStmt()
	execParams := dynamicModelFieldConditionValues(s.whereModel, s.selectWhereFields)

	rows, err := s.stmt.Query(execParams...)
	if err != nil {
		mysqlQueryErrorLog(err)
		return nil, err
	}
	defer rows.Close()

	var results []interface{}
	modelValue := reflect.ValueOf(s.whereModel)
	for rows.Next() {
		tempModel := reflect.New(modelValue.Type()).Elem()
		var selectFieldsAddr []interface{}

		for _, field := range s.selectFields {
			// 将字段转为帕斯卡写法
			tempFieldName := strings.Replace(field, "_", " ", -1)
			tempFieldName = strings.Title(tempFieldName)
			tempFieldName = strings.Replace(tempFieldName, " ", "", -1)
			selectFieldsAddr = append(selectFieldsAddr, tempModel.FieldByName(tempFieldName).Addr().Interface())
		}

		if len(s.selectFields) == 0 {
			for i := 0; i < tempModel.NumField(); i++ {
				selectFieldsAddr = append(selectFieldsAddr, tempModel.Field(i).Addr().Interface())
			}
		}

		err := rows.Scan(selectFieldsAddr...)
		if err != nil {
			return nil, err
		}
		results = append(results, tempModel.Interface())
	}
	return results, nil
}
