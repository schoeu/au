package analysis

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	tagsUrlArr    = map[string][]string{}
	tagsMaxLength int
	tagsLimit     bool
)

func TagsUrl(filePath string, mLength int, lmt bool) {
	tagsMaxLength = mLength
	tagsLimit = lmt
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
			getTags(tmpStr)
		}
	}
}

func getTags(c string) {
	tagsInfo := pluginRe.FindAllStringSubmatch(c, -1)
	if len(tagsInfo) > 0 && len(tagsInfo[0]) > 1 {
		url := tagsInfo[0][1]
		tags := tagsInfo[0][2]
		tagsArr := strings.Split(tags, ", ")

		for _, v := range tagsArr {
			item := tagsUrlArr[v]
			if (len(item) < tagsMaxLength) || !tagsLimit {
				tagsUrlArr[v] = append(item, url)
			}
		}
	}
}

func GetTagsMap(cwd string) map[string][]string {
	dir := ensureDir(cwd)

	b, err := json.Marshal(tagsUrlArr)
	if err != nil {
		log.Fatal(err)
	}

	finalPath := filepath.Join(dir, "tags_urls"+tempExt)
	fmt.Printf("\nTags file in %v\n", finalPath)
	if e := ioutil.WriteFile(finalPath, b, 0777); e != nil {
		log.Fatal(e)
	}
	return tagsUrlArr
}
