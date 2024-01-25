package config

import (
	"errors"
	"strings"

	"github.com/fimreal/goutils/ezap"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.StringP("port", "p", "6788", "serve address")
	pflag.BoolP("debug", "d", false, "enable debug")

	pflag.StringP("provider", "P", "tencentcloud", "specify dns provider")
	pflag.String("secretid", "", `provider's secret id`)
	pflag.String("secretkey", "", `provider's secret key`)

	pflag.String("ddns", "", "give a domain to enable ddns")
	pflag.String("ddns.spec", "*/10 * * * *", "crond schedule")

	//
	pflag.ErrHelp = errors.New("")
	// shell 不允许带'.'的环境变量，识别环境变量时去除'.'
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, ``))
	viper.AutomaticEnv()
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	if viper.GetBool("debug") {
		ezap.SetLevel("debug")
	}
}
