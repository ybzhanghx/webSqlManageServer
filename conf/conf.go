package conf

import (
	"flag"

	"github.com/BurntSushi/toml"
)

type (
	LoggerConfig struct {
		FileName            string `json:"filename"`
		Level               int    `json:"level"`    // 日志保存的时候的级别，默认是 Trace 级别
		Maxlines            int    `json:"maxlines"` // 每个文件保存的最大行数，默认值 1000000
		Maxsize             int    `json:"maxsize"`  // 每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
		Daily               bool   `json:"daily"`    // 是否按照每天 logrotate，默认是 true
		Maxdays             int    `json:"maxdays"`  // 文件最多保存多少天，默认保存 7 天
		Rotate              bool   `json:"rotate"`   // 是否开启 logrotate，默认是 true
		Perm                string `json:"perm"`     // 日志文件权限
		RotatePerm          string `json:"rotateperm"`
		EnableFuncCallDepth bool   `json:"-"` // 输出文件名和行号
		LogFuncCallDepth    int    `json:"-"` // 函数调用层级
	}
	LogConf struct {
		LogPath string
		Level   int
		Adapter string
	}
	CommConf struct {
		Ver     string
		LogConf *LogConf
	}
	Config struct {
		*CommConf
		Server    *Server
		MysqlConf *MysqlConf
		FuncList  *FuncListConf
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
		Types    string
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
	flag.StringVar(&confPath, "conf", "./WebManageSvr.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
