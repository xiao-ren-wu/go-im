package setup

import (
	"github.com/xiao-ren-wu/go-im/configs"
	"github.com/xiao-ren-wu/go-im/middleware"
)

func Init() {
	configs.InitConf()
	middleware.InitLoggo(configs.App.Loggo)
}
