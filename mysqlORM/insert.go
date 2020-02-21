package mysqlORM

import (
	"log"
	"reflect"
	"strings"
)

type Insert struct {
	opBasic
	insertFields []string
}


func (i *Insert)SetInsertFields(f []string) *Insert {
	i.setStmtNil()
	i.insertFields = f
	i.insertFields = append(i.insertFields, "created_at", "updated_at")
	return i
}

func (i *Insert)prepareStmt(){
	if i.stmt == nil {
		if len(i.insertFields) > 0 {
			var insertPlaceholder []string
			for count := 0; count < len(i.insertFields); count++ {
				insertPlaceholder = append(insertPlaceholder, "?")
			}
			sqlInsert := "insert into " + i.table +
				" (" + strings.Join(i.insertFields, ",") + ") values (" + strings.Join(insertPlaceholder, ",") + ")"

			stmt, err := i.db.Prepare(sqlInsert)

			if err != nil {
				log.Print("insert prepare failed:"+err.Error())
			} else {
				i.stmt = stmt
			}
		} else {
			log.Print("without insert fields been set")
		}
	}
}

func (i *Insert)Commit(model interface{}) (interface{}, error) {
	i.prepareStmt()

	m := reflect.New(reflect.TypeOf(model)).Elem()
	m.Set(reflect.ValueOf(model))

	dataTimeNow := mysqlDateTimeNow()
	m.FieldByName("CreatedAt").SetString(dataTimeNow)
	m.FieldByName("UpdatedAt").SetString(dataTimeNow)

	insertValues := dynamicModelFieldValues(m.Interface(), i.insertFields)

	res, err := i.stmt.Exec(insertValues...)
	if err != nil {
		mysqlQueryErrorLog(err)
		return model, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		mysqlQueryErrorLog(err)
		return model, err
	}

	m.FieldByName("Id").SetInt(lastId)

	return m.Interface(), nil
}