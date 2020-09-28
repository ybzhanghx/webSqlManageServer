package test

import (
	"WebManageSvr/conf"
	"WebManageSvr/controllers"
	"WebManageSvr/mysqls"
	_ "WebManageSvr/routers"
	"WebManageSvr/service"
	"bytes"
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

func TestGetTableList2(t *testing.T) {
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}

	mysqls.MysqlInit()

	_, _ = service.GetDBNames()

}

func TestMe(t *testing.T) {
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}

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

func TestGetDB(t *testing.T) {
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}

	mysqls.MysqlInit()
	var firstB []byte
	for i := 0; i < 10; i++ {
		d, _ := service.GetDBNames()

		b, _ := json.Marshal(d)
		if i == 0 {
			firstB = b
		} else {
			if !bytes.Equal(firstB, b) {
				fmt.Println(string(firstB))
				fmt.Println(string(b))
				fmt.Println(d)
				fmt.Println(i)
				break
			}
		}

	}

}
func TestGos(t *testing.T) {
	a := map[string]int{"2": 2, "3": 3, "4": 4}
	for i := range a {
		fmt.Println(a[i])
	}
}
