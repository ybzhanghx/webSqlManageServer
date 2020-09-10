package conf

import (
	"flag"

	"bailun.com/CT4_quote_server/common/conf"
	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		*conf.CommConf
		Server    *Server
		MysqlConf *MysqlConf
		funcList  *FuncListConf
	}
	Server struct {
		AppName         string
		HTTPPort        int
		RunMode         string
		CopyRequestBody bool
		AutoRender      bool
		EnableDocs      bool
	}

	MysqlConf struct {
		Addr         string
		User         string
		Pwd          string
		DbName       string
		SystemDbName string
	}

	FuncListConf struct {
		Name     string
		Value    string
		types    string
		Children []FuncNode
	}
	FuncNode struct {
		Name     string
		Value    string
		Children []FuncNode
	}
)

var (
	confPath string
	Conf     *Config
)

func init() {
	flag.StringVar(&confPath, "conf", "./conf.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
