package service

import (
	"bailun.com/CT4_quote_server/WebManageSvr/models"
	"bailun.com/CT4_quote_server/WebManageSvr/mysqls"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
)

//获取 数据库表数据
func GetDBNames() (res []models.DBTBInfo, err error) {
	sqlFmt := `select TABLE_SCHEMA DB,table_name TB  from information_schema.tables
			where  table_type='BASE TABLE' ORDER BY TABLE_SCHEMA;`

	var sqlData []models.DbTb
	err = mysqls.SysInfDb.Select(&sqlData, sqlFmt)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	var tmpMap = make(map[string][]string)
	for i := range sqlData {
		if _, ok := tmpMap[sqlData[i].DB]; !ok {
			tmpMap[sqlData[i].DB] = []string{}
		}
		tmpMap[sqlData[i].DB] = append(tmpMap[sqlData[i].DB], sqlData[i].TB)
	}

	for key, value := range tmpMap {
		var tmpNode models.DBTBInfo
		tmpNode.DbName = key
		tmpNode.TbName = value
		res = append(res, tmpNode)
	}
	return
}

func GetDBTbConfig(dbtb *models.DbTb) (res []models.DataTableConfigReturn, err error) {
	var sqlData []models.DataTableConfig
	if sqlData, err = models.ReadDBTBConfig(dbtb); err != nil {
		return
	}

	err = copier.Copy(&res, &sqlData)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	for i, v := range sqlData {
		res[i].IsKey = v.ColumnKey == "PRI"
		res[i].IsAbleNull = v.IsNull == "YES"
	}
	return
}

//更新配置
func UpdateDBConfig(db, tb string, data []models.DataTableUpdateConfig) (err error) {

	var dbs *sqlx.DB

	var ok bool
	if dbs, ok = mysqls.ArrSqlDb[db]; !ok {
		return errors.New("not found db")
	}

	tx, err := dbs.Begin()

	var funcFind = func(name string) (res []int) {
		for i := range data {
			if data[i].Action == name {
				res = append(res, i)
			}
		}
		return
	}
	addList := funcFind("add")
	updateList := funcFind("update")
	delList := funcFind("del")

	addsql := "alter table " + tb + " add column ("
	for i, v := range addList {
		if i != 0 {
			addsql += ", "
		}
		addsql += fmt.Sprintf("%s %s(%d)", data[v].NewName, data[v].DataType, data[v].Length)

	}
	addsql += `)`
	if len(addList) != 0 {
		_, err = tx.Exec(addsql)
		if err != nil {
			mysqls.RoBackMysqlFunc(err, tx)
			return
		}
	}

	updatesql := "alter table " + tb + " "
	for i, v := range updateList {
		if i != 0 {
			updatesql += ", "
		}
		updatesql += fmt.Sprintf("change  %s %s %s(%d) ", data[v].FieldName, data[v].NewName, data[v].DataType,
			data[v].Length)
	}
	if len(updateList) != 0 {
		_, err = tx.Exec(updatesql)
		if err != nil {
			mysqls.RoBackMysqlFunc(err, tx)
			return
		}
	}

	delsql := "alter table " + tb + "  "
	for i, v := range delList {
		if i != 0 {
			addsql += ", "
		}
		delsql += fmt.Sprintf("drop  %s ", data[v].FieldName)
	}

	if len(delList) != 0 {
		_, err = tx.Exec(delsql)
	}
	if err != nil {
		mysqls.RoBackMysqlFunc(err, tx)
		return
	}
	err = tx.Commit()
	if err != nil {
		logs.Error(err.Error())
		return
	}

	return
}
