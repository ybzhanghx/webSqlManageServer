package models

import (
	"bailun.com/CT4_quote_server/WebManageSvr/mysqls"
	"bailun.com/CT4_quote_server/WebManageSvr/utils"
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
	WHERE  LOWER(table_schema) = "%s" and TABLE_NAME = "%s"`, dbtb.DB, dbtb.TB)
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
