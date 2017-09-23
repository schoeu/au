package analysis

import (
	"net/url"
	"log"
	"strings"
)

var domainSuffix = []string{
	".com",".la",".io",".co",".info",".net",".org",".me",".mobi",
	".us",".biz",".xxx",".ca",".co.jp",".com.cn",".net.cn",".edu.cn",
	".org.cn",".mx",".tv",".ws",".ag",".com.ag",".net.ag", ".cn",
	".org.ag",".am",".asia",".at",".be",".com.br",".net.br",
	".bz",".com.bz",".net.bz",".cc",".com.co",".net.co",
	".nom.co",".de",".es",".com.es",".nom.es",".org.es",
	".eu",".fm",".fr",".gs",".in",".co.in",".firm.in",".gen.in",
	".ind.in",".net.in",".org.in",".it",".jobs",".jp",".ms",
	".com.mx",".nl",".nu",".co.nz",".net.nz",".org.nz", "design",
	".se",".tc",".tk",".tw",".com.tw",".idv.tw",".org.tw",
	".hk",".co.uk",".me.uk",".org.uk",".vg", ".com.hk"}

func GetDomain(l string) string{

	u, err := url.Parse(l)
	if err != nil {
		log.Fatal(err)
	}

	host := u.Host

	for _, v := range domainSuffix {
		if strings.Contains(host, v) {
			tempStr := strings.Replace(host, v, "", -1)
			splitRs := strings.Split(tempStr, ".")
			top := splitRs[len(splitRs)-1:]
			return top[0] + v
		}
	}
	return host
}
