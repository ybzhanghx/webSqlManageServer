package conf

import (
	"bailun.com/CT4_quote_server/lib/sqltool"
	"github.com/jmoiron/sqlx"
)

var (
	SysInfDb *sqlx.DB
	ArrSqlDb map[string]*sqlx.DB
)

func MysqlInit() {
	c := Conf.MysqlConf
	var err error

	SysInfDb, err = sqltool.InitDB(c.User, c.Pwd, c.Addr, c.SystemDbName)
	if err != nil {
		panic(err)
	}
	ArrSqlDb = make(map[string]*sqlx.DB)
	var dbName []struct {
		DataBase string `db:"Database"`
	}

	err = SysInfDb.Select(&dbName, "show databases ")
	if err != nil {
		panic(err)
	}

	for i := range dbName {
		var tmpDb *sqlx.DB
		tmpDb, err = sqltool.InitDB(c.User, c.Pwd, c.Addr, dbName[i].DataBase)
		if err != nil {
			panic(err)
		}
		ArrSqlDb[dbName[i].DataBase] = tmpDb
	}

	//testInfDb, err = sqltool.InitDB(conf.Conf.MysqlConf.User, conf.Conf.MysqlConf.Pwd, conf.Conf.MysqlConf.Addr,
	//	"zybtest")
	//tradeFxDb, err = sqltool.InitDB(conf.Conf.MysqlConf.User, conf.Conf.MysqlConf.Pwd, conf.Conf.MysqlConf.Addr,
	//	"TradeFxDB")
	//if err != nil {
	//	panic(err)
	//}
	//funcTableMap.Store("test", dataBaseMapTable{"zybtest", "tableA"})
	//funcTableMap.Store("clientManage", dataBaseMapTable{"CT4DB", "info_user"})
	//dbNames, _ := GetDBNames("TradeFxDB")
	//for i := range dbNames {
	//	funcTableMap.Store(dbNames[i], dataBaseMapTable{"TradeFxDB", dbNames[i]})
	//}
}
