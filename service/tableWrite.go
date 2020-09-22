package service

import (
	"WebManageSvr/models"
	"WebManageSvr/mysqls"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	"github.com/ompluscator/dynamic-struct"
)

func WriteTable(dbInfo *models.DbTb, add, upd string, del []int) (err error) {
	dbs, err2 := mysqls.GetDbs(dbInfo.DB)
	if err2 != nil {
		return err2
	}
	sqlFmt := fmt.Sprintf("SELECT * FROM `%s` limit 1", dbInfo.TB)
	var rows *sqlx.Rows
	rows, err = dbs.Queryx(sqlFmt)
	if err != nil {
		logs.Error(err)
		return err
	}
	//dbs.Select()

	colField, _ := rows.ColumnTypes()

	typeStruct := dynamicstruct.NewStruct()
	colTypes := make([]models.FieldType, len(colField))
	for i := range colField {
		colTypes[i] = models.TypeDatabaseNameNoNull(typeStruct, colField[i])
	}
	pData := typeStruct.Build().NewSliceOfStructs()

	tx, err := dbs.Begin()

	if len(add) != 0 {
		if err = models.TableInsert(dbInfo.TB, tx, add, pData, colTypes); err != nil {
			return
		}
	}
	if len(del) != 0 {
		if err = models.DelTableRows(dbInfo.TB, tx, del); err != nil {
			return
		}
	}
	if len(upd) != 0 {
		if err = models.UpdateTableRows(dbInfo.TB, tx, upd, pData, colTypes); err != nil {
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	return nil
}
