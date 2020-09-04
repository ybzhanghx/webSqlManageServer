package models

import (
	"bailun.com/CT4_quote_server/WebManageSvr/conf"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"log"
	"strconv"
	"strings"
)

type WriteParm struct {
	DB  string
	TB  string
	Add string
	Del []int
	Upd string
}
type tt interface {
}

func WriteTable(p WriteParm) (err error) {
	var dbs *sqlx.DB
	var ok bool
	if dbs, ok = conf.ArrSqlDb[p.DB]; !ok {
		return errors.New("not found db")
	}
	sqlFmt := fmt.Sprintf("SELECT * FROM `%s` ", p.TB)
	var rows *sqlx.Rows
	rows, err = dbs.Queryx(sqlFmt)
	if err != nil {
		logs.Error(err)
		return err
	}
	//dbs.Select()

	colField, _ := rows.ColumnTypes()

	typeStruct := dynamicstruct.NewStruct()
	colTypes := make([]FieldType, len(colField))
	for i := range colField {
		colTypes[i] = typeDatabaseNameNoNull(typeStruct, colField[i])
	}
	pData := typeStruct.Build().NewSliceOfStructs()

	tx, err := dbs.Begin()
	if len(p.Add) != 0 {
		if err = Add(p.TB, tx, p.Add, pData, colTypes); err != nil {
			return
		}
	}
	if len(p.Del) != 0 {
		if err = DelTableRows(p.TB, tx, p.Del); err != nil {
			return
		}
	}
	if len(p.Upd) != 0 {
		if err = UpdateTableRows(p.TB, tx, p.Upd, pData, colTypes); err != nil {
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
func roBackFunc(err error, tx *sql.Tx) {
	logs.Error(err.Error())
	err = tx.Rollback()
	if err != nil {
		log.Println("tx.Rollback() Error:" + err.Error())
		return
	}

}
func Add(tb string, tx *sql.Tx, bytes string, pData interface{}, colTypes []FieldType) (err error) {
	getBytes := []byte(bytes)
	err = json.Unmarshal(getBytes, &pData)
	if err != nil {
		fmt.Println(err.Error())
	}
	reader := dynamicstruct.NewReader(pData)
	readersSlice := reader.ToSliceOfReaders()
	getFieldName := func() []string {
		var slices []string
		for i, v := range colTypes {
			if strings.ToLower(v.FieldName) != "id" {
				slices = append(slices, colTypes[i].FieldName)
			}
		}
		return slices
	}()
	getValueFunc := func(readersSliceNode dynamicstruct.Reader) []string {
		var slices []string
		for _, v := range colTypes {
			if strings.ToLower(v.FieldName) == "id" {
				continue
			}
			structFieldName := strings.ToUpper(v.FieldName[0:1]) + v.FieldName[1:] //首字母大写
			getV := readersSliceNode.GetField(structFieldName)
			var tmpV string
			switch v.TypeName {
			case "int":
				tmpV = strconv.Itoa(getV.Int())
			case "string":
				tmpV = "'" + getV.String() + "'"
			case "time":
				tmpV = "'" + getV.String() + "'"
			}
			slices = append(slices, tmpV)
		}
		return slices
	}

	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) values", tb, strings.Join(getFieldName, ","))
	var arrValue []string
	for i := range readersSlice {
		itemTmp := "(" + strings.Join(getValueFunc(readersSlice[i]), ",") + ")"
		arrValue = append(arrValue, itemTmp)
	}
	sqlStr += strings.Join(arrValue, ",")

	_, err = tx.Exec(sqlStr)
	if err != nil {
		roBackFunc(err, tx)
		return
	}
	return
}

func UpdateTableRows(tb string, tx *sql.Tx, bytes string, pData interface{}, colTypes []FieldType) (err error) {
	_ = json.Unmarshal([]byte(bytes), &pData)
	reader := dynamicstruct.NewReader(pData)
	readersSlice := reader.ToSliceOfReaders()

	getValueFunc := func(readersSliceNode dynamicstruct.Reader) (idValue string, upV []string) {

		for _, v := range colTypes {

			structFieldName := strings.ToUpper(v.FieldName[0:1]) + v.FieldName[1:] //首字母大写
			if structFieldName == "Id" {
				continue
			}
			getV := readersSliceNode.GetField(structFieldName)
			var tmpV string
			switch v.TypeName {
			case "int":
				tmpV = strconv.Itoa(getV.Int())
			case "string":
				tmpV = "'" + getV.String() + "'"
			case "time":
				tmpV = "'" + getV.String() + "'"
			}
			setV := v.FieldName + "=" + tmpV
			upV = append(upV, setV)
		}

		idValue = strconv.Itoa(readersSliceNode.GetField("Id").Int())
		return
	}

	for i := range readersSlice {
		idv, slicev := getValueFunc(readersSlice[i])
		sqlStr := fmt.Sprintf("update %s SET %s WHERE id = %s ", tb, strings.Join(slicev, ","), idv)
		_, err = tx.Exec(sqlStr)
		if err != nil {
			roBackFunc(err, tx)
			return
		}

	}

	return
}

func DelTableRows(tb string, tx *sql.Tx, value []int) (err error) {

	delIdArr := func() []string {
		var tmp []string
		for i := range value {
			tmp = append(tmp, strconv.Itoa(value[i]))
		}
		return tmp
	}()
	sqlFmt := fmt.Sprintf("DELETE FROM %s WHERE id in (%s)", tb, strings.Join(delIdArr, ","))
	_, err = tx.Exec(sqlFmt)
	if err != nil {
		roBackFunc(err, tx)
		return
	}
	return

}
