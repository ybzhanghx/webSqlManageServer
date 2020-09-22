package models

import (
	"WebManageSvr/mysqls"
	"WebManageSvr/utils"
	"database/sql"
	"fmt"
	"github.com/astaxie/beego/logs"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"strings"
)

type (
	DbTb struct {
		DB string `db:"DB"`
		TB string `db:"TB"`
	}
	DBTBInfo struct {
		DbName string
		TbName []string
	}
	DataTableConfig struct {
		FieldName string `db:"field_name" json:"field_name"`
		FieldDesc string `db:"field_desc" json:"field_desc"`
		DataType  string `db:"data_type" json:"data_type"`
		IsNull    string `db:"is_null" json:"is_null"`
		Length    int    `db:"length" json:"length"`
		ColumnKey string `db:"COLUMN_KEY" json:"columnKey"`
	}
	DataTableConfigReturn struct {
		FieldName  string `json:"field_name"`
		FieldDesc  string `json:"field_desc"`
		DataType   string `json:"data_type"`
		IsNull     string `json:"is_null"`
		Length     int    `json:"length"`
		IsKey      bool   `json:"is_key"`
		IsAbleNull bool   `json:"is_able_null"`
	}

	DataTableUpdateConfig struct {
		DataTableConfig `json:",inline"`
		NewName         string `json:"newName"`
		Action          string `json:"action"`
	}

	FieldType struct {
		FieldName string
		TypeName  string
		AbleNull  bool
	}
)

func ReadDBTBConfig(dbtb *DbTb) (sqlData []DataTableConfig, err error) {
	sqlFmt := fmt.Sprintf(`SELECT COLUMN_NAME field_name,column_comment field_desc,DATA_TYPE data_type,
	IS_NULLABLE is_null,IFNULL(CHARACTER_MAXIMUM_LENGTH,0) length,COLUMN_KEY 
	FROM INFORMATION_SCHEMA.COLUMNS
	WHERE  TABLE_SCHEMA = "%s" and TABLE_NAME = "%s"`, dbtb.DB, dbtb.TB)
	if err = mysqls.SysInfDb.Select(&sqlData,
		sqlFmt); err != nil {
		logs.Error(sqlData, err)
		return
	}
	return
}

//获取表类型数据
func TypeDatabaseName(newStruct dynamicstruct.Builder, Field *sql.ColumnType) (colTypes FieldType) {
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

}

//获取表类型数据
func TypeDatabaseNameNoNull(newStruct dynamicstruct.Builder, Field *sql.ColumnType) (colTypes FieldType) {
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
			return 0
		case strings.Contains(typeName, "CHAR"), strings.Contains(typeName, "TEXT"):
			colTypes.TypeName = "string"

			return ""
		case strings.Contains(typeName, "DATETIME"):
			colTypes.TypeName = "time"

			return ""
		default:
			return 0
		}
	}()

	newStruct = newStruct.AddField(newTypeName, types, tagStr)
	return

}

func GetAllDBTB(isSystem bool) (sqlData []DbTb, err error) {
	sqlFmt := `select TABLE_SCHEMA DB,table_name TB  from information_schema.tables
			where  table_type='BASE TABLE' ORDER BY DB;`
	err = mysqls.SysInfDb.Select(&sqlData, sqlFmt)
	if err != nil {
		logs.Error(err.Error())
	}

	if !isSystem {
		j := 0
		for _, val := range sqlData {
			var tmp = strings.ToLower(val.DB)
			if tmp != "information_schema" && tmp != "mysql" && tmp != "performance_schema" {
				sqlData[j] = val
				j++
			}
		}
		sqlData = sqlData[:j]
	}
	//sort.Sort()
	return
}
