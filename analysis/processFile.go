package analysis

import (
	"../autils"
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

var (
	re         = regexp.MustCompile("http[s]?://\\w+\\S*\\b")
	ignorExts  = [4]string{".jpg", ".png", ".gif", ".jpeg"}
	uniqUrlMap = map[string]int{}
	fileName   = ""
	tempExt    = ".atmp"
)

type rsMapType map[string][]string

type siteInfo struct {
	Top   string   `json:"top"`
	Sites []string `json:"sites"`
}

type siteCtt []siteInfo

// 日志处理入口
func Process(path string, cwd string, name string) {
	fileName = name
	// 读取文件
	readLine(path)
	makeMap(cwd)
}

// 单行读取日志
func readLine(filePath string) {
	uniqUrlMap = map[string]int{}
	fi, err := os.Open(filePath)
	autils.ErrHadle(err)
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		content := string(a)
		analysisFile(content)
	}
}

// 信息保存到strict中
func makeMap(cwd string) {
	rsMap := rsMapType{}
	for k, _ := range uniqUrlMap {
		top := autils.GetDomain(k)
		host := top.Host
		if host != "" {
			rsMap[host] = append(rsMap[host], k)
		}
	}
	MergeInfos(cwd, rsMap)
}

// 日志分析
func analysisFile(content string) {
	rs := re.FindAllStringSubmatch(content, -1)
	url := rs[0][0]
	crtExt := filepath.Ext(url)
	for _, v := range ignorExts {
		if v == crtExt {
			return
		}
	}
	matched, _ := regexp.MatchString("^http[s]?:\\/\\/.+?", url)
	if matched {
		uniqUrlMap[url] = 1
	}
}
