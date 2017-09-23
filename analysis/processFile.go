package analysis

import (
	"os"
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"regexp"
)

var (
	re         = regexp.MustCompile("http[s]?://\\w+\\S*\\b")
	ignorExts  = [...]string{".jpg", ".png", ".gif", ".jpeg"}
	ctt        = []string{}
	uniqUrlMap = map[string]int{}
	rsMap      = map[string][]string{}
	fileName   = ""
	tempRs     = "result"
)

func Process(path string, cwd string, name string) {

	fileName = name
	// 读取文件
	readLine(path)

	for k, _ := range uniqUrlMap {
		ctt = append(ctt, k)
	}
	makeMap(cwd)
}

// 单行读取日志
func readLine(filePath string) {
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
	for k, _ := range uniqUrlMap {
		top := GetDomain(k)

		rsMap[top] = append(rsMap[top], k)
	}

	dir := ensureDir(cwd)

	for k, v := range rsMap {
		tmpfn := filepath.Join(dir, k)
		content := strings.Join(v, "\n")
		if err := ioutil.WriteFile(tmpfn, []byte(content), 0777); err != nil {
			log.Fatal(err)
		}
	}
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