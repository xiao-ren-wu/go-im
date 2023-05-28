package configs

import (
	"embed"
	"fmt"
	"github.com/xiao-ren-wu/go-im/model"
	"github.com/xiao-ren-wu/gonfig"
	"gopkg.in/yaml.v3"
)

//go:embed *.yaml
var confRaw embed.FS

var App model.AppConf

func InitConf() {
	var err error
	if err = gonfig.Unmarshal(confRaw, &App, gonfig.FilePrefix("conf")); err != nil {
		panic(err)
	}
	raw, _ := yaml.Marshal(App)
	fmt.Printf("read conf: \n%s\n", string(raw))
}
