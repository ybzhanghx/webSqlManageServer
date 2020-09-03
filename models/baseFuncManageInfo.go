package models

import (
	"bailun.com/CT4_quote_server/WebManageSvr/conf"
	"bailun.com/CT4_quote_server/WebManageSvr/utils"

	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"strings"

	"database/sql"
	"errors"
	"github.com/astaxie/beego/logs"
	"log"

	"fmt"
	"github.com/jinzhu/copier"
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
	ColumnKey string `db:"COLUMN_KEY" json:"columnKey"`
}
type DataTableConfigReturn struct {
	FieldName  string `json:"field_name"`
	FieldDesc  string `json:"field_desc"`
	DataType   string `json:"data_type"`
	IsNull     string `json:"is_null"`
	Length     int    `json:"length"`
	IsKey      bool   `json:"is_key"`
	IsAbleNull bool   `json:"is_able_null"`
}

type DataTableUpdateConfig struct {
	DataTableConfig `json:",inline"`
	NewName         string `json:"newName"`
	Action          string `json:"action"`
}

func GetDBTbConfig(db, tb string) (res []DataTableConfigReturn, err error) {
	sqlfmt := fmt.Sprintf(`SELECT COLUMN_NAME field_name,column_comment field_desc,DATA_TYPE data_type,
	IS_NULLABLE is_null,IFNULL(CHARACTER_MAXIMUM_LENGTH,0) length,COLUMN_KEY 
	FROM INFORMATION_SCHEMA.COLUMNS
	WHERE  LOWER(table_schema) = "%s" and TABLE_NAME = "%s"`, db, tb)
	var sqlData []DataTableConfig
	if err = conf.SysInfDb.Select(&sqlData,
		sqlfmt); err != nil {

		logs.Error(sqlData, err)
		return
	}
	//res = make([]DataTableConfigReturn,len(sqlData))
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

//获取 数据库表数据
func GetDBNames() (res []DBTBInfo, err error) {
	sqlFmt := `select TABLE_SCHEMA DB,table_name TB  from information_schema.tables
			where  table_type='BASE TABLE' ORDER BY TABLE_SCHEMA;`

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

//获取表数据
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

//更新配置
func UpdateDBConfig(db, tb string, data []DataTableUpdateConfig) (err error) {

	var dbs *sqlx.DB

	var ok bool
	if dbs, ok = conf.ArrSqlDb[db]; !ok {
		return errors.New("not found db")
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
	err = tx.Commit()
	if err != nil {
		logs.Error(err.Error())
		return
	}

	return
}

//获取表数据
func GetTableDataList(db, tb string, page, size int) (data []interface{}, colTypes []FieldType, err error) {
	var dbs *sqlx.DB
	var ok bool
	if dbs, ok = conf.ArrSqlDb[db]; !ok {
		return nil, nil, errors.New("not found db")
	}

	sqlFmt := fmt.Sprintf("SELECT * FROM `%s`  Limit %d,%d", tb, (page-1)*size, size)
	var rows *sqlx.Rows
	rows, err = dbs.Queryx(sqlFmt)
	if err != nil {
		logs.Error(err)
		return nil, nil, err
	}
	//dbs.Select()

	colField, _ := rows.ColumnTypes()

	typeStruc := dynamicstruct.NewStruct()
	colTypes = make([]FieldType, len(colField))
	for i := range colField {
		colTypes[i] = typeDatabaseName(typeStruc, colField[i])
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

//获取表产犊
func GetTableDataTotals(db, tb string) (res int64, err error) {
	var dbs *sqlx.DB
	var ok bool
	if dbs, ok = conf.ArrSqlDb[db]; !ok {
		return 0, errors.New("not found db")
	}
	sqlFmt := fmt.Sprintf("SELECT count(*) FROM `%s` ", tb)
	err = dbs.QueryRow(sqlFmt).Scan(&res)
	if err != nil {
		logs.Error(err.Error())
	}
	return
}

type FieldType struct {
	FieldName string
	TypeName  string
	AbleNull  bool
}

//获取表类型数据
func typeDatabaseName(newStruct dynamicstruct.Builder, Field *sql.ColumnType) (colTypes FieldType) {
	FieldName := Field.Name()
	typeName := Field.DatabaseTypeName()
	isNull, _ := Field.Nullable()

	newTypeName := strings.ToUpper(FieldName[0:1]) + FieldName[1:]
	tagStr := fmt.Sprintf(`db:"%s" json:"%s"`, FieldName, FieldName)

	colTypes.AbleNull = isNull
	colTypes.FieldName = FieldName
	var types interface{}
	types = func() interface{} {
		switch {
		case strings.Contains(typeName, "INT"):
			colTypes.TypeName = "int"
			if !isNull {
				return 0
			} else {
				return sql.NullInt64{}
			}
		case strings.Contains(typeName, "CHAR"), strings.Contains(typeName, "TEXT"):
			colTypes.TypeName = "string"
			if !isNull {
				return ""
			} else {
				return sql.NullString{}
			}

		case strings.Contains(typeName, "DATETIME"):
			colTypes.TypeName = "time"
			return utils.NullTimeStamp{}
		default:
			return 0
		}
	}()

	newStruct = newStruct.AddField(newTypeName, types, tagStr)
	return
	//switch{
	//case strings.Contains(typeName,"INT"):
	//
	//	//types = 0
	//	newStruct = newStruct.AddField(newTypeName,0, `db:"`+FieldName+`"`)
	//case strings.Contains(typeName,"CHAR"),strings.Contains(typeName,"TEXT"):
	//	newStruct = newStruct.AddField(newTypeName,"", `db:"`+FieldName+`"`)
	//case strings.Contains(typeName,"DATETIME"):
	//	newStruct = newStruct.AddField(newTypeName,"", `db:"`+FieldName+`"`)
	//}
	//AddField("DataBase", "", `json:"dataBasse" db:"Database"`).
	//	Build()
	//switch mf.fieldType {
	//case fieldTypeBit:
	//	return "BIT"
	//case fieldTypeBLOB:
	//	if mf.charSet != collations[binaryCollation] {
	//		return "TEXT"
	//	}
	//	return "BLOB"
	//case fieldTypeDate:
	//	return "DATE"
	//case fieldTypeDateTime:
	//	return "DATETIME"
	//case fieldTypeDecimal:
	//	return "DECIMAL"
	//case fieldTypeDouble:
	//	return "DOUBLE"
	//case fieldTypeEnum:
	//	return "ENUM"
	//case fieldTypeFloat:
	//	return "FLOAT"
	//case fieldTypeGeometry:
	//	return "GEOMETRY"
	//case fieldTypeInt24:
	//	return "MEDIUMINT"
	//case fieldTypeJSON:
	//	return "JSON"
	//case fieldTypeLong:
	//	return "INT"
	//case fieldTypeLongBLOB:
	//	if mf.charSet != collations[binaryCollation] {
	//		return "LONGTEXT"
	//	}
	//	return "LONGBLOB"
	//case fieldTypeLongLong:
	//	return "BIGINT"
	//case fieldTypeMediumBLOB:
	//	if mf.charSet != collations[binaryCollation] {
	//		return "MEDIUMTEXT"
	//	}
	//	return "MEDIUMBLOB"
	//case fieldTypeNewDate:
	//	return "DATE"
	//case fieldTypeNewDecimal:
	//	return "DECIMAL"
	//case fieldTypeNULL:
	//	return "NULL"
	//case fieldTypeSet:
	//	return "SET"
	//case fieldTypeShort:
	//	return "SMALLINT"
	//case fieldTypeString:
	//	if mf.charSet == collations[binaryCollation] {
	//		return "BINARY"
	//	}
	//	return "CHAR"
	//case fieldTypeTime:
	//	return "TIME"
	//case fieldTypeTimestamp:
	//	return "TIMESTAMP"
	//case fieldTypeTiny:
	//	return "TINYINT"
	//case fieldTypeTinyBLOB:
	//	if mf.charSet != collations[binaryCollation] {
	//		return "TINYTEXT"
	//	}
	//	return "TINYBLOB"
	//case fieldTypeVarChar:
	//	if mf.charSet == collations[binaryCollation] {
	//		return "VARBINARY"
	//	}
	//	return "VARCHAR"
	//case fieldTypeVarString:
	//	if mf.charSet == collations[binaryCollation] {
	//		return "VARBINARY"
	//	}
	//	return "VARCHAR"
	//case fieldTypeYear:
	//	return "YEAR"
	//default:
	//	return ""
	//}
}

//删除表数据
func DelTableRow(db, tb string, value []string) (err error) {
	var dbs *sqlx.DB
	var ok bool
	if dbs, ok = conf.ArrSqlDb[db]; !ok {
		err = errors.New("not found db")
		logs.Error(err.Error())
		return err
	}

	rofunc := func(err error, tx *sql.Tx) {
		logs.Error(err.Error())
		err = tx.Rollback()
		if err != nil {
			log.Println("tx.Rollback() Error:" + err.Error())
			return
		}

	}
	var tx *sql.Tx
	tx, err = dbs.Begin()
	delIdArr := func() string {
		tmp := ""
		for i := range value {
			tmp += value[i] + ","
		}
		return tmp[:len(tmp)-1]
	}
	sqlFmt := fmt.Sprintf("DELETE FROM %s WHERE id in (%s）", tb, delIdArr)
	_, err = tx.Exec(sqlFmt)
	if err != nil {
		rofunc(err, tx)
		return
	}

	if err = tx.Commit(); err != nil {
		logs.Error(err.Error())
		return
	}
	return

}
