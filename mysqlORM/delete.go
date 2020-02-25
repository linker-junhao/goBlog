package mysqlORM

import "log"

type Delete struct {
	opBasic
	deleteWhereFields whereCondition
	whereModel        interface{}
}

func (d *Delete) SetDeleteWhereFields(calculateType string, fields []string) *Delete {
	d.setStmtNil()
	d.deleteWhereFields[calculateType] = fields
	return d
}

func (d *Delete) Where(model interface{}) *Delete {
	d.whereModel = model
	return d
}

func (d *Delete) prepareStmt() {
	if d.stmt == nil {
		theSql := "delete from " + d.table + makeConditionClause(d.deleteWhereFields)
		stmt, err := d.db.Prepare(theSql)
		if err != nil {
			log.Print("delete prepare failed:" + err.Error())
		}
		d.stmt = stmt
	}
}

func (d *Delete) Commit() (interface{}, error) {
	d.prepareStmt()
	execParams := dynamicModelFieldConditionValues(d.whereModel, d.deleteWhereFields)
	_, err := d.stmt.Exec(execParams...)
	if err != nil {
		mysqlQueryErrorLog(err)
	}
	return d.whereModel, nil
}
