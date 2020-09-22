package service

import (
	"WebManageSvr/models"
	"WebManageSvr/mysqls"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/iancoleman/orderedmap"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
)

//获取 数据库表数据
func GetDBNames() (res []models.DBTBInfo, err error) {
	var sqlData []models.DbTb
	sqlData, err = models.GetAllDBTB(false)
	if err != nil {
		return
	}

	var tmpMap = orderedmap.New()
	for i := range sqlData {

		var resValue []string
		if res, ok := tmpMap.Get(sqlData[i].DB); ok {
			resValue = res.([]string)
		}
		resValue = append(resValue, sqlData[i].TB)
		tmpMap.Set(sqlData[i].DB, resValue)
	}

	for _, key := range tmpMap.Keys() {
		var tmpNode models.DBTBInfo
		tmpNode.DbName = key
		tmpV, _ := tmpMap.Get(key)
		tmpNode.TbName = tmpV.([]string)
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
