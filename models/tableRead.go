package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
)

func ReadTableRowsByPage(tb string, page int, size int, dbs *sqlx.DB) (rows *sqlx.Rows, err error) {
	sqlFmt := fmt.Sprintf("SELECT * FROM `%s`  Limit %d,%d", tb, (page-1)*size, size)
	rows, err = dbs.Queryx(sqlFmt)
	if err != nil {
		logs.Error(err)
	}
	return
}

//获取表产犊
func ReadTableDataTotals(tb string, dbs *sqlx.DB) (res int64, err error) {
	sqlFmt := fmt.Sprintf("SELECT count(*) FROM `%s` ", tb)
	err = dbs.QueryRow(sqlFmt).Scan(&res)
	if err != nil {
		logs.Error(err.Error())
	}
	return
}
