package conf

import (
	"flag"

	"bailun.com/CT4_quote_server/common/conf"
	"github.com/BurntSushi/toml"
)

type Server struct {
	AppName         string
	HTTPPort        int
	RunMode         string
	CopyRequestBody bool
	AutoRender      bool
	EnableDocs      bool
}

type MysqlConf struct {
	Addr         string
	User         string
	Pwd          string
	DbName       string
	SystemDbName string
}

type Config struct {
	*conf.CommConf
	Server    *Server
	MysqlConf *MysqlConf
}

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
