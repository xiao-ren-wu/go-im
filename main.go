package main

import (
	"github.com/xiao-ren-wu/go-im/middleware"
	"github.com/xiao-ren-wu/go-im/setup"
)

func main() {
	setup.Init()
	middleware.L.Info("server start success")
}
