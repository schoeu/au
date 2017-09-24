package analysis

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type tType map[string]int

var (
	tagRe    = regexp.MustCompile("\\[mip-tags used\\]")
	pluginRe = regexp.MustCompile("\\[mip-tags used\\](http[s]?://\\S+): ([\\s\\S]*) log queue")
	tagsMap  = tType{}
)

// 单行读取日志
func CountData(filePath string) {
	fi, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		tmpStr := string(a)
		if tagRe.MatchString(tmpStr) {
			analyTags(tmpStr)
			break
		}

	}
}

func GetCountData() tType {
	return tagsMap
}

func analyTags(c string) {
	tagsInfo := pluginRe.FindAllStringSubmatch(c, -1)
	if len(tagsInfo) > 0 && len(tagsInfo[0]) > 1 {
		tags := tagsInfo[0][2]
		tagArr := strings.Split(tags, ", ")
		for _, v := range tagArr {
			tagsMap[v] += 1
		}
	}
}
