package ddns

import (
	"strings"
	"time"

	"github.com/fimreal/goutils/ezap"
	"github.com/fimreal/goutils/http"
	"github.com/fimreal/watson/component/adapter"
	"github.com/fimreal/watson/component/crond"
	"github.com/fimreal/watson/component/utils"
	"github.com/fimreal/watson/config"
	"github.com/fimreal/watson/providers"
)

func Hold(domain string) {
	if domain == "" {
		return
	}
	ezap.Info("开启定时更新域名解析动态地址 [" + domain + "]")
	p, err := adapter.NewProvider()
	if err != nil {
		ezap.Fatal(err.Error())
	}
	zone, r, err := utils.ParseDomain(domain)
	if err != nil {
		ezap.Fatal(err.Error())
	}
	r.Type = "A"

	ezap.Info("初始化...")
	if err = DDNS(p, zone, r); err != nil {
		ezap.Fatal(err)
	}
	crond.Run("*/10 * * * *", func() {
		_ = (DDNS(p, zone, r))
	})
}

func DDNS(p providers.DNSProvider, zone string, r providers.Record) error {
	ip, err := http.HttpGet(config.WhatIsMyIP)
	if err != nil {
		return err
	}
	r.TTL = time.Duration(600)
	r.Value = strings.Trim(ip, "\n")
	return p.SetRecord(zone, r)
}
