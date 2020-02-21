package mysqlORM

import (
	"log"
	"reflect"
	"strings"
)

type Update struct {
	opBasic
	updateFields    []string
	updateWhereFields whereCondition
	whereModel interface{}
}


func (upd *Update)SetUpdateFields(f []string) *Update {
	upd.setStmtNil()
	upd.updateFields = f
	upd.updateFields = append(upd.updateFields, "updated_at")
	return upd
}

func (upd *Update)SetUpdateWhereFields(calculateType string, fields []string) *Update {
	upd.setStmtNil()
	upd.updateWhereFields[calculateType] = fields
	return upd
}

func (upd *Update)Where(model interface{}) *Update {
	upd.whereModel = model
	return upd
}

func (upd *Update)prepareStmt(){
	if upd.stmt == nil {
		sqlUpdate := "Update " + upd.table + " set " + strings.Join(upd.updateFields, " = ?, ") + " = ? where " + makeConditionClause(upd.updateWhereFields)
		stmt, err := upd.db.Prepare(sqlUpdate)
		if err != nil {
			log.Print("Update prepare failed:"+err.Error())
		}
		upd.stmt = stmt
	}
}

func (upd *Update)Commit(model interface{}) (interface{}, error) {
	upd.prepareStmt()

	retM := reflect.New(reflect.TypeOf(model)).Elem()
	retM.Set(reflect.ValueOf(model))
	retM.FieldByName("UpdatedAt").SetString(mysqlDateTimeNow())

	execParams := append(dynamicModelFieldValues(retM.Interface(), upd.updateFields), dynamicModelFieldConditionValues(upd.whereModel, upd.updateWhereFields)...)
	_, err := upd.stmt.Exec(execParams...)
	if err != nil {
		mysqlQueryErrorLog(err)
		return model, err
	}
	return retM.Interface(), nil
}