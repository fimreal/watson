package adapter

import (
	"fmt"

	"github.com/fimreal/watson/providers"
	"github.com/fimreal/watson/providers/tencentcloud"
	"github.com/spf13/viper"
)

func NewProvider() (providers.DNSProvider, error) {
	provider := viper.GetString("provider")
	switch provider {
	// case "alidns":
	// 	return alidns.NewProvider()
	case "tencentcloud":
		return tencentcloud.NewProvider()
	default:
		return nil, fmt.Errorf("unrecognized DNS Provider: %s", provider)
	}
}
