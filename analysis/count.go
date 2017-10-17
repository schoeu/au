package analysis

import (
	"../autils"
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type tType map[string]int

var (
	tagRe    = regexp.MustCompile("\\[mip-tags used\\]")
	pluginRe = regexp.MustCompile("\\[mip-tags used\\](http[s]?://\\S+): ([\\s\\S]*) log queue")
	tagsMap  = tType{}
	relReg   = regexp.MustCompile(`["|']`)
)

// 单行读取日志
func CountData(filePath string) {
	fi, err := os.Open(filePath)
	autils.ErrHadle(err)
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		tmpStr := string(a)

		if !relReg.MatchString(tmpStr) && tagRe.MatchString(tmpStr) {
			analyTags(tmpStr)
		}

	}
}

func GetCountData(cwd string, anaDate string) {
	var bf bytes.Buffer
	bf.WriteString("UPDATE tags SET url_count = CASE tag_name")
	for k, v := range tagsMap {
		bf.WriteString(" WHEN '")
		bf.WriteString(k)
		bf.WriteString("' THEN ")
		bf.WriteString(strconv.Itoa(v))

	}
	bf.WriteString(" END ")
	bf.WriteString(` WHERE tags.ana_date = '`)
	bf.WriteString(anaDate)
	bf.WriteString(`'`)

	sqlStr := bf.String()
	db := autils.OpenDb(cwd)

	/*
		UPDATE categories
		SET display_order = CASE id
		WHEN 1 THEN 3
		WHEN 2 THEN 4
		WHEN 3 THEN 5
		END
		WHERE id IN (1,2,3)
	*/

	_, err := db.Exec(sqlStr)
	autils.ErrHadle(err)
	defer db.Close()
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
