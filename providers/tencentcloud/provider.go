package tencentcloud

import (
	"errors"

	"github.com/spf13/viper"
)

func NewProvider() (*Provider, error) {
	p := &Provider{
		SecretId:  viper.GetString("secretid"),
		SecretKey: viper.GetString("secretkey"),
	}

	if p.SecretId == "" || p.SecretKey == "" {
		return p, errors.New("没找到 SecretId 或者 SecretKey")
	}
	return p, nil
}
