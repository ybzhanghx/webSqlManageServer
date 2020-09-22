package mysqls

import (
	"WebManageSvr/conf"
	"database/sql"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	SysInfDb *sqlx.DB
	ArrSqlDb map[string]*sqlx.DB
)

func MysqlInit() {
	c := conf.Conf.MysqlConf
	err := ConnectNewSql(c)
	if err != nil {
		panic(err)
	}
}
func ConnectNewSql(c *conf.MysqlConf) (err error) {

	SysInfDb, err = InitDB(c.User, c.Pwd, c.Addr, "INFORMATION_SCHEMA")
	if err != nil {
		logs.Error(err.Error())
		return
	}
	ArrSqlDb = make(map[string]*sqlx.DB)
	var dbName []struct {
		DataBase string `db:"Database"`
	}
	err = SysInfDb.Select(&dbName, "show databases ")
	if err != nil {
		logs.Error(err.Error())
	}

	for i := range dbName {
		var tmpDb *sqlx.DB
		tmpDb, err = InitDB(c.User, c.Pwd, c.Addr, dbName[i].DataBase)
		if err != nil {
			logs.Error(err.Error())
		}
		ArrSqlDb[dbName[i].DataBase] = tmpDb
	}
	return
}

func GetDbs(DB string) (*sqlx.DB, error) {
	var dbs *sqlx.DB
	var ok bool
	if dbs, ok = ArrSqlDb[DB]; !ok {
		return nil, errors.New("not found db")
	}
	return dbs, nil
}

func RoBackMysqlFunc(err error, tx *sql.Tx) {
	logs.Error(err.Error())
	err = tx.Rollback()
	if err != nil {
		logs.Error("tx.Rollback() Error:" + err.Error())
		return
	}
}

func InitDB(dbuser, dbpwd, dbhost, dbname string) (db *sqlx.DB, err error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbuser, dbpwd, dbhost, dbname)

	db, err = sqlx.Connect("mysql", dns)
	if err != nil {
		logs.Error("db connect failed", err)
		return
	}
	// defer DB.Close()
	logs.Info("db connect succeed")
	return
}
