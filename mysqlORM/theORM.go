package mysqlORM

import (
	"database/sql"
	"log"
)

type fields []string

type whereCondition map[string]fields

type MysqlORM struct {
	config ORMConfig
}

func NewMysqlORM(config ORMConfig) MysqlORM {
	return MysqlORM{
		config: config,
	}
}

func (orm *MysqlORM)SetConfig(config ORMConfig) MysqlORM {
	orm.config = config
	return *orm
}

func (orm MysqlORM)Delete() *Delete {
	return &Delete{
		opBasic: opBasic{
			db:                orm.config.db,
			table:             orm.config.table,
			stmt:              nil,
		},
		deleteWhereFields: whereCondition{},
		whereModel:        nil,
	}
}


func (orm MysqlORM)Select() *Select {
	return &Select{
		opBasic: opBasic{
			db:                orm.config.db,
			table:             orm.config.table,
			stmt:              nil,
		},
		selectFields:      []string{},
		selectWhereFields: whereCondition{},
		whereModel:        nil,
	}
}


func (orm MysqlORM)Update() *Update {
	return &Update{
		opBasic: opBasic{
			db:                orm.config.db,
			table:             orm.config.table,
			stmt:              nil,
		},
		updateFields:      []string{},
		updateWhereFields: whereCondition{},
		whereModel:        nil,
	}
}

func (orm MysqlORM)Insert() *Insert {
	return &Insert{
		opBasic: opBasic{
			db:                orm.config.db,
			table:             orm.config.table,
			stmt:              nil,
		},
		insertFields:      []string{},
	}
}


type opBasic struct {
	db *sql.DB
	table string
	stmt *sql.Stmt
}


func (o *opBasic)setStmtNil(){
	if o.stmt != nil {
		err := o.stmt.Close()
		if err != nil {
			log.Println("stmt close err: ",err)
		}
		o.stmt = nil
	}
}
