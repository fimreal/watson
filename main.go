package main

import (
	"github.com/fimreal/watson/component/ddns"
	"github.com/fimreal/watson/serve"
	"github.com/spf13/viper"
)

func main() {
	ddns.Hold(viper.GetString("ddns"))
	serve.Run(":" + viper.GetString("port"))
}
