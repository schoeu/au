package analysis

import (
	"../autils"
	"../config"
	"bufio"
	"bytes"
	"database/sql"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type tagsUrlType map[string][]string
type rsType struct {
	list  []string
	count int
}

var (
	tagsUrlArr   = tagsUrlType{}
	tagsRsUrlArr = tagsUrlType{}
	tagRsPath    string
	tagRelReg    = regexp.MustCompile(`["|']`)
)

// 保存数据只保留前十条
const tagMax = 10

// 单小时组件信息处理入口
func TagsUrl(filePath string, cwd string, fileName string) {
	tagsUrlArr = tagsUrlType{}
	tagsRsUrlArr = tagsUrlType{}
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
		if !tagRelReg.MatchString(tmpStr) && tagRe.MatchString(tmpStr) {
			getTags(tmpStr)
		}
	}

	var buf bytes.Buffer
	for k, v := range tagsUrlArr {
		// 过滤空tag
		if k == "" {
			continue
		}

		buf.WriteString(k)
		buf.WriteString(" ")

		b, uDArr := getDiffUrls(v)
		l := len(v)
		if l > tagMax {
			l = tagMax
		}

		buf.WriteString(strings.Join(b[:l], ","))
		buf.WriteString(" ")

		for key, val := range uDArr {
			buf.WriteString("," + key + "=" + strconv.Itoa(val))
		}
		buf.WriteString("\n")
	}

	tagRsPath = autils.EnsureDir(filepath.Join(cwd, config.TagTempDir))
	finalPath := filepath.Join(tagRsPath, fileName+tempExt)
	err = ioutil.WriteFile(finalPath, []byte(buf.String()), 0777)
	autils.ErrHadle(err)
}

// 组件信息去重
func getDiffUrls(val []string) ([]string, map[string]int) {
	var uniqUrlArr, normalArr []string
	uniqDomainArr := map[string]int{}
	for _, v := range val {
		d := autils.GetDomain(v).Host
		if uniqDomainArr[d] == 0 {
			uniqUrlArr = append(uniqUrlArr, v)
		} else {
			normalArr = append(normalArr, v)
		}
		uniqDomainArr[d] += 1
	}
	a := append(uniqUrlArr, normalArr...)
	return a, uniqDomainArr
}

// 获取组件信息
func getTags(c string) {
	tagsInfo := pluginRe.FindAllStringSubmatch(c, -1)
	if len(tagsInfo) > 0 && len(tagsInfo[0]) > 1 {
		url := tagsInfo[0][1]
		tags := tagsInfo[0][2]
		tagsArr := strings.Split(tags, ", ")

		for _, v := range tagsArr {
			tagsUrlArr[v] = append(tagsUrlArr[v], url)
		}
	}
}

// 组件组件信息，写入数据库
func GetTagsMap(anaDate string, db *sql.DB) {
	tagCountCtt := map[string]string{}
	files, err := ioutil.ReadDir(tagRsPath)
	autils.ErrHadle(err)
	for _, file := range files {
		fi, err := os.Open(filepath.Join(tagRsPath, file.Name()))
		autils.ErrHadle(err)
		defer fi.Close()
		br := bufio.NewReader(fi)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			infos := bytes.Split(a, []byte(" "))
			if len(infos) > 2 {
				tag := string(infos[0])
				urlArr := strings.Split(string(infos[1]), ",")
				tagC := infos[2][1:]
				tagCountCtt[tag] = string(tagC)
				if len(tagsRsUrlArr[tag]) > 0 {
					tagsRsUrlArr[tag] = append(tagsRsUrlArr[tag], urlArr...)
				} else {
					tagsRsUrlArr[tag] = urlArr
				}
			}
		}
	}

	for k, v := range tagsRsUrlArr {
		sort.Strings(v)
		tagsRsUrlArr[k] = uniq(v)
	}
	bArr := []string{}
	for k, v := range tagsRsUrlArr {
		rl := len(v)
		if rl > tagMax {
			rl = tagMax
		}

		rs, _ := getDiffUrls(v)
		tmp := strings.Join(rs[:rl], ",")

		tagCountStr := tagCountCtt[k]
		tagCountNum := strings.Split(tagCountStr, ",")
		bArr = append(bArr, "('"+k+"', '"+tmp+"', '0', '"+string(tagCountStr)+"','"+strconv.Itoa(len(tagCountNum))+"','"+anaDate+"', '"+autils.GetCurrentData(time.Now())+"')")
	}

	autils.ErrHadle(err)

	sqlStr := "INSERT INTO tags (tag_name, urls, url_count, tag_count, domain_count, ana_date, edit_date) VALUES " + strings.Join(bArr, ",")
	_, err = db.Exec(sqlStr)
	autils.ErrHadle(err)
}
