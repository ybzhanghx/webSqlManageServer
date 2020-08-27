package models

import (
	"bailun.com/CT4_quote_server/WebManageSvr/conf"

	"database/sql"
	"errors"
	"github.com/astaxie/beego/logs"
	"log"

	"fmt"

	"github.com/jmoiron/sqlx"
	"sync"
)

const ()

type dataBaseMapTable struct {
	dataBase string
	table    string
}

var (
	funcTableMap sync.Map
)

func BaseFuncManageInit() {
	//var err error
	//sysInfDb, err = sqltool.InitDB(conf.Conf.MysqlConf.User, conf.Conf.MysqlConf.Pwd, conf.Conf.MysqlConf.Addr,
	//	conf.Conf.MysqlConf.SystemDbName)
	//
	//testInfDb, err = sqltool.InitDB(conf.Conf.MysqlConf.User, conf.Conf.MysqlConf.Pwd, conf.Conf.MysqlConf.Addr,
	//	"zybtest")
	//tradeFxDb, err = sqltool.InitDB(conf.Conf.MysqlConf.User, conf.Conf.MysqlConf.Pwd, conf.Conf.MysqlConf.Addr,
	//	"TradeFxDB")
	//if err != nil {
	//	panic(err)
	//}
	funcTableMap.Store("test", dataBaseMapTable{"zybtest", "tableA"})
	funcTableMap.Store("clientManage", dataBaseMapTable{"CT4DB", "info_user"})
	dbNames, _ := GetTBNamesByDB("TradeFxDB")
	for i := range dbNames {
		funcTableMap.Store(dbNames[i], dataBaseMapTable{"TradeFxDB", dbNames[i]})
	}
}

func GetDBTBInfo(funcName string) (db, tb string, ok bool) {
	var tmp interface{}
	tmp, ok = funcTableMap.Load(funcName)
	if !ok {
		return
	}
	tbdb := tmp.(dataBaseMapTable)
	db, tb = tbdb.dataBase, tbdb.table
	return
}

type DataTableConfig struct {
	FieldName string `db:"field_name" json:"field_name"`
	FieldDesc string `db:"field_desc" json:"field_desc"`
	DataType  string `db:"data_type" json:"data_type"`
	IsNull    string `db:"is_null" json:"is_null"`
	Length    int    `db:"length" json:"length"`
}

type DataTableUpdateConfig struct {
	DataTableConfig `json:",inline"`
	NewName         string `json:"newName"`
	Action          string `json:"action"`
}

func GetDBTbConfig(db, tb string) (res []DataTableConfig, err error) {
	sqlfmt := fmt.Sprintf(`SELECT COLUMN_NAME field_name,column_comment field_desc,DATA_TYPE data_type,
	IS_NULLABLE is_null,IFNULL(CHARACTER_MAXIMUM_LENGTH,0) length
	FROM INFORMATION_SCHEMA.COLUMNS
	WHERE  LOWER(table_schema) = "%s" and TABLE_NAME = "%s"`, db, tb)

	if err = conf.SysInfDb.Select(&res,
		sqlfmt); err != nil {

		logs.Error(res, err)
		return
	}
	return
}

//type DBTBName struct {
//	DbName string
//	TbName []string
//}
type DBTBInfo struct {
	DbName string
	TbName []string
}
type tmpSqlData struct {
	DB string `db:"DB"`
	TB string `db:"TB"`
}

func GetDBNames() (res []DBTBInfo, err error) {
	sqlFmt := `select TABLE_SCHEMA DB,table_name TB  from information_schema.tables
			where  table_type='base table' ORDER BY TABLE_SCHEMA;`

	var sqlData []tmpSqlData
	err = conf.SysInfDb.Select(&sqlData, sqlFmt)
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
		var tmpNode DBTBInfo
		tmpNode.DbName = key
		tmpNode.TbName = value
		res = append(res, tmpNode)
	}
	return
}

func GetTBNamesByDB(db string) (res []string, err error) {
	sqlfmt := fmt.Sprintf(`select table_name from information_schema.tables 	
		where table_schema='%s' and table_type='base table'`, db)
	if db == "" {
		sqlfmt = `select table_name from information_schema.tables
			where table_schema=%s and table_type='base table`
	}

	if err = conf.SysInfDb.Select(&res,
		sqlfmt); err != nil {
		//if err == sql.ErrNoRows {
		//	err = ErrUserIsNotExist
		//}

		logs.Error(res, err)
		return
	}
	return
}

func UpdateDBConfig(db, tb string, data []DataTableUpdateConfig) (err error) {

	var dbs *sqlx.DB

	if db == "zybtest" {
		dbs = conf.SysInfDb
	} else {
		return
	}
	rofunc := func(err error, tx *sql.Tx) {
		logs.Error(err.Error())
		err = tx.Rollback()
		if err != nil {
			log.Println("tx.Rollback() Error:" + err.Error())
			return
		}

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
			rofunc(err, tx)
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
			rofunc(err, tx)
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
		rofunc(err, tx)
		return
	}

	return
}

func GetTableDataList(db, tb string, page, size int) (data []map[string]interface{}, err error) {
	var dbs *sqlx.DB
	if db == "zybtest" {
		dbs = conf.SysInfDb
	} else if db == "TradeFxDB" {
		dbs = conf.SysInfDb
	} else {
		return nil, errors.New("not found db")
	}
	//data =make([]interface{},0)

	sqlFmt := fmt.Sprintf("SELECT * FROM `%s`  Limit %d,%d", tb, (page-1)*size, size)
	var rows *sqlx.Rows
	rows, err = dbs.Queryx(sqlFmt)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	for rows.Next() {
		//下面演示如何将数据保存到struct、map和数组中
		//定义struct对象
		//var p Place

		//定义map类型
		m := make(map[string]interface{})

		////定义slice类型
		//s := make([]interface{}, 0)
		//
		////使用StructScan函数将当前记录的数据保存到struct对象中
		//err = rows.StructScan(&p)
		////保存到map
		err = rows.MapScan(m)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		data = append(data, m)
		//保存到数组
		//err = rows.SliceScan(&s)
	}
	//if err = dbs.Select(&data, sqlFmt); err != nil {
	//	if err == sql.ErrNoRows {
	//		err = ErrUserIsNotExist
	//	}
	//
	//	return
	//}
	fmt.Println(data)
	return

}
