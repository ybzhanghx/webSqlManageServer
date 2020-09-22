package service

import (
	"WebManageSvr/models"
	"WebManageSvr/mysqls"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	"github.com/ompluscator/dynamic-struct"
)

//获取表数据
func GetTableDataList(dbtb *models.DbTb, page, size int) (data []interface{}, colTypes []models.FieldType, err error) {
	var dbs *sqlx.DB
	if dbs, err = mysqls.GetDbs(dbtb.DB); err != nil {
		return
	}

	var rows *sqlx.Rows
	if rows, err = models.ReadTableRowsByPage(dbtb.TB, page, size, dbs); err != nil {
		return
	}
	//dbs.Select()

	colField, _ := rows.ColumnTypes()
	typeStruc := dynamicstruct.NewStruct()
	colTypes = make([]models.FieldType, len(colField))
	for i := range colField {
		colTypes[i] = models.TypeDatabaseName(typeStruc, colField[i])
	}
	buildStruct := typeStruc.Build()

	for rows.Next() {
		var node = buildStruct.New()
		err = rows.StructScan(node)
		if err != nil {
			logs.Error(err)
		}
		err = nil
		data = append(data, node)
	}
	return
}

//获取表总数
func GetTableDataTotals(dbtb *models.DbTb) (res int64, err error) {
	var dbs *sqlx.DB
	if dbs, err = mysqls.GetDbs(dbtb.DB); err != nil {
		return
	}
	return models.ReadTableDataTotals(dbtb.TB, dbs)
}
