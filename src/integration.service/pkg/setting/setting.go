package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg          *ini.File
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize   int
	JwtSecrect string
)

// init 函数
func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	// 如果初始化配置文件失败 那么报错消息
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app,ini' : %v", err)
	}
	LoadBase()
	LoadServer()
	LoadApp()

}

func LoadBase() {
	RunMode = Cfg.Section("").Key("run_mode").MustString("debug")
}

// 服务基础加载
func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}
	HTTPPort = sec.Key("http_port").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("read_timeout").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("write_timeout").MustInt(60)) * time.Second
}

// 加载app json web token密码分页等
func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecrect = sec.Key("jwt_secret").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("page_size").MustInt(20)
}
