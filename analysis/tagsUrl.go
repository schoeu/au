package analysis

import (
	"os"
	"log"
	"bufio"
	"io"
	"strings"
)

type tagsType struct {
	name string
	list []string
}

var (
	tagsUrlArr = []tagsType{}
)

func TagsUrl(filePath string) {
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
			break
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
			if len(tagsUrlArr) == 0 {
				s := []string{url}
				t := tagsType{ v, s}
				tagsUrlArr = append(tagsUrlArr, t)
				continue
			}
			for i, val := range tagsUrlArr {
				if v == val.name {
					item := tagsUrlArr[i].list
					if (len(item) < maxLength) || !limit {
						tagsUrlArr[i].list = append(item, url)
					}
				} else {
					s := []string{url}
					t := tagsType{ v, s}
					tagsUrlArr = append(tagsUrlArr, t)
				}
			}

		}
	}
}

func GetTagsMap() []tagsType{
	return tagsUrlArr
}