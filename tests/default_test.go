package test

import (
	"bailun.com/CT4_quote_server/WebManageSvr/conf"
	"bailun.com/CT4_quote_server/WebManageSvr/controllers"
	"bailun.com/CT4_quote_server/WebManageSvr/mysqls"
	_ "bailun.com/CT4_quote_server/WebManageSvr/routers"
	"bailun.com/CT4_quote_server/WebManageSvr/service"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jmoiron/sqlx"
	"github.com/ompluscator/dynamic-struct"
	"path/filepath"
	"runtime"
	"testing"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestGetTableList(t *testing.T) {
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	conf.BeeConfInit()
	mysqls.MysqlInit()

	_, _, _ = service.GetTableDataList("TradeFxDB", "trade_account", 1, 10)
	//func (mf *mysqlField) typeDatabaseName() string {
	//	switch mf.fieldType {
	//case fieldTypeBit:
	//	return "BIT"
	//case fieldTypeBLOB:
	//	if mf.charSet != collations[binaryCollation] {
	//	return "TEXT"
	//}
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
	//	return "LONGTEXT"
	//}
	//	return "LONGBLOB"
	//case fieldTypeLongLong:
	//	return "BIGINT"
	//case fieldTypeMediumBLOB:
	//	if mf.charSet != collations[binaryCollation] {
	//	return "MEDIUMTEXT"
	//}
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
	//	return "BINARY"
	//}
	//	return "CHAR"
	//case fieldTypeTime:
	//	return "TIME"
	//case fieldTypeTimestamp:
	//	return "TIMESTAMP"
	//case fieldTypeTiny:
	//	return "TINYINT"
	//case fieldTypeTinyBLOB:
	//	if mf.charSet != collations[binaryCollation] {
	//	return "TINYTEXT"
	//}
	//	return "TINYBLOB"
	//case fieldTypeVarChar:
	//	if mf.charSet == collations[binaryCollation] {
	//	return "VARBINARY"
	//}
	//	return "VARCHAR"
	//case fieldTypeVarString:
	//	if mf.charSet == collations[binaryCollation] {
	//	return "VARBINARY"
	//}
	//	return "VARCHAR"
	//case fieldTypeYear:
	//	return "YEAR"
	//default:
	//	return ""
	//}
	//}
}

func TestGetTableList2(t *testing.T) {
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	conf.BeeConfInit()
	mysqls.MysqlInit()

	_, _ = service.GetDBNames()

}

func TestMe(t *testing.T) {
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	conf.BeeConfInit()
	mysqls.MysqlInit()
	typeStruc := dynamicstruct.NewStruct().
		AddField("DataBase", "", `json:"dataBasse" db:"Database"`).
		Build()
	getData := typeStruc.New()

	var rows *sqlx.Rows
	//err = conf.SysInfDb.Select(&getData, "show databases ")
	rows, err = mysqls.SysInfDb.Queryx("show databases ")
	if err != nil {
		fmt.Println(err.Error())
	}
	for rows.Next() {
		err = rows.StructScan(getData)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(getData)
	}

	fmt.Println(getData)
}

func TestTime(t *testing.T) {
	bytes := []byte(
		`{"TableInsert":"[{\"id\":\"row_43\",\"itemA\":\"ddsg\",\"itemTime\":\"2020-09-04 00:01:00\"}]","Del":[],"Upd":"",
"DB":"zybtest","TB":"tableA"}`)
	data := controllers.UpdateTableParm{}
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func TestInsert(t *testing.T) {

}
