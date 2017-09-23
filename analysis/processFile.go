package analysis

import (
	"bufio"
	"io"
	"os"
	//"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	//"encoding/json"
	"strings"
)

var (
	re         = regexp.MustCompile("http[s]?://\\w+\\S*\\b")
	ignorExts  = [...]string{".jpg", ".png", ".gif", ".jpeg"}
	uniqUrlMap = map[string]int{}
	fileName   = ""
	tempRs     = "result"
	tempExt    = ".atmp"
	// 一个站点最多保存多个少url
	maxLength = 10
	notlimit  = false
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
	if err != nil {
		log.Fatal(err)
		return
	}
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

// 创建临时文件夹存放中间文件
func ensureDir(cwd string) string {
	dirPath := filepath.Join(cwd, tempRs)
	mkDirErr := os.MkdirAll(dirPath, 0777)
	if mkDirErr != nil {
		log.Fatal(mkDirErr)
	}
	return dirPath
}

// 信息保存到strict中
func makeMap(cwd string) {
	rsMap := rsMapType{}
	for k, _ := range uniqUrlMap {
		top := GetDomain(k)
		host := top.host
		scheme := top.scheme
		replacedUrl := strings.Replace(k, host, "*", -1)
		replacedUrl = strings.Replace(replacedUrl, scheme+"://", "", 1)
		key := scheme + "@" + host
		if (len(rsMap[key]) < maxLength) || notlimit {
			rsMap[key] = append(rsMap[key], replacedUrl)
		}
	}
	MergeInfos(rsMap)
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

	uniqUrlMap[url] = 1
}
