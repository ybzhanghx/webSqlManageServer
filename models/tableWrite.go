package models

import (
	"WebManageSvr/mysqls"
	"database/sql"
	"encoding/json"
	"fmt"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"strconv"
	"strings"
)

func TableInsert(tb string, tx *sql.Tx, bytes string, pData interface{}, colTypes []FieldType) (err error) {
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
		mysqls.RoBackMysqlFunc(err, tx)
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
				fmt.Println(getV.Int())
				tmpV = strconv.Itoa(getV.Int())
			case "string":
				tmpV = "'" + getV.String() + "'"
			case "time":
				tmpV = "'" + getV.String() + "'"
			}
			setV := "`" + v.FieldName + "` =" + tmpV
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
			mysqls.RoBackMysqlFunc(err, tx)
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
		mysqls.RoBackMysqlFunc(err, tx)
		return
	}
	return

}
