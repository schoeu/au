package analysis

import (
	"os"
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"encoding/json"
	"fmt"
	"strings"
)

var (
	re         = regexp.MustCompile("http[s]?://\\w+\\S*\\b")
	ignorExts  = [...]string{".jpg", ".png", ".gif", ".jpeg"}
	uniqUrlMap = map[string]int{}
	fileName   = ""
	tempRs     = "result"
)

type rsMapType map[string][]string

type siteInfo struct {
	Top    string `json:"top"`
	Sites   []string `json:"sites"`
}

type siteCtt [] siteInfo

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

func ensureDir(cwd string) string{
	dirPath := filepath.Join(cwd, tempRs, fileName)
	mkDirErr := os.MkdirAll(dirPath, 0777)
	if mkDirErr != nil {
		log.Fatal(mkDirErr)
	}
	return dirPath
}

func makeMap(cwd string) {
	rsMap := rsMapType{}
	for k, _ := range uniqUrlMap {
		top := GetDomain(k)
		host := top.host
		scheme := top.scheme
		replacedUrl := strings.Replace(k, host, "*",-1)
		replacedUrl = strings.Replace(replacedUrl, scheme + "://", "",1)
		key := host + "@" + scheme
		rsMap[key] = append(rsMap[key], replacedUrl)
	}

	dir := ensureDir(cwd)

	b, err := json.Marshal(rsMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	if err := ioutil.WriteFile(dir + ".jox", b, 0777); err != nil {
		log.Fatal(err)
	}

	//for k, v := range rsMap {
	//	tmpfn := filepath.Join(dir, k)
	//	content := strings.Join(v, "\n")
	//	if err := ioutil.WriteFile(tmpfn, []byte(content), 0777); err != nil {
	//		log.Fatal(err)
	//	}
	//}
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