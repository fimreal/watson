package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/fimreal/watson/providers"
)

// 分割域名改为 dns 记录格式，例如
// 1. www.domain.com => www domain.com
// 2. domain.com => @ domain.com
// 3. domain => error()
func SplitDomain(domain string) (sd string, zone string, err error) {
	re := regexp.MustCompile(`[^.]+\.[^.]+$`)
	zone = re.FindString(domain)
	if zone == "" {
		return "", "", errors.New("format domain err: " + domain)
	}
	sd = strings.TrimSuffix(domain, "."+zone)
	if sd == zone {
		sd = "@"
	}
	return
}

func ParseDomain(domain string) (zone string, record providers.Record, err error) {
	sd, zone, err := SplitDomain(domain)
	if err != nil {
		return
	}
	record.Name = sd
	return
}
