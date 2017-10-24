package analysis

import (
	"log"
	"net/url"
	"strings"
)

var domainSuffix = []string{
	".com.ag", ".org.tw", ".org.nz", ".org.uk",
	".co.jp", ".edu.cn", ".com.bz", ".com.co", ".net.co",
	".com.br", ".org.ag", ".org.cn", ".net.cn", ".com.cn", "design",
	".nom.co", ".com.es", ".nom.es", ".co.in", ".firm.in", ".gen.in",
	".co.uk", ".me.uk", ".com.hk", ".com.tw", ".idv.tw",
	".hk", ".se", ".tc", ".tk", ".tw", ".vg", ".jobs", ".jp", ".de",
	".cn", ".com", ".la", ".io", ".co", ".info", ".net", ".org", ".me",
	".mobi", ".us", ".biz", ".xxx", ".ca", ".mx", ".tv", ".ws", ".ag", ".cc", ".bz",
	".asia", ".at", ".be", ".eu", ".it",
}

type urlInfos struct {
	host   string
	scheme string
}

// 获取域名主域
func GetDomain(l string) urlInfos {

	uInfo := urlInfos{}

	u, err := url.Parse(l)
	if err != nil {
		log.Fatal(err)
	}

	host := u.Host
	uInfo.host = host
	uInfo.scheme = u.Scheme

	for _, v := range domainSuffix {
		if strings.Contains(host, v) {
			tempStr := strings.Replace(host, v, "", -1)
			splitRs := strings.Split(tempStr, ".")
			top := splitRs[len(splitRs)-1:]
			uInfo.host = top[0] + v
			return uInfo
		}
	}
	return uInfo
}
