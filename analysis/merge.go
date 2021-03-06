package analysis

import (
	"../autils"
	"../config"
	"bufio"
	"bytes"
	"database/sql"
	"io"
	"io/ioutil"
	"log"
	//"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type mit map[string]int

var (
	rsPath string
	m      = mit{}
)

type uniqInfoType map[string][]string

// 域名数据处理入口
func MergeInfos(cwd string, msg rsMapType) {
	var bf bytes.Buffer
	m = mit{}
	for k, v := range msg {
		l := len(v)
		m[k] = m[k] + l
		if l > 10 {
			l = 10
		}

		bf.WriteString(k)
		bf.WriteString(" ")
		bf.WriteString(strings.Join(v[:l], ","))
		bf.WriteString(" ")
		bf.WriteString(strconv.Itoa(l))
		bf.WriteString("\n")
	}

	rsPath = autils.EnsureDir(filepath.Join(cwd, config.TempDir))
	finalPath := filepath.Join(rsPath, fileName+tempExt)
	if e := ioutil.WriteFile(finalPath, []byte(bf.String()), 0777); e != nil {
		log.Fatal(e)
	}
}

// 域名数据写入数据库
func CalcuUniqInfo(anaDate string, db *sql.DB) {
	t := uniqInfoType{}
	files, err := ioutil.ReadDir(rsPath)
	autils.ErrHadle(err)

	for _, file := range files {
		fi, err := os.Open(filepath.Join(rsPath, file.Name()))
		autils.ErrHadle(err)
		defer fi.Close()
		br := bufio.NewReader(fi)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			content := string(a)
			infos := strings.Split(content, " ")
			if len(infos) > 1 {
				tmpUrl := infos[1]
				if tmpUrl != "" {
					//u, err := url.QueryUnescape(tmpUrl)
					//if err == nil {
					//	tmpUrl = u
					//}

					infos[1] = strings.Replace(tmpUrl, "'", "''", -1)
				}
				tag := infos[0]
				urlArr := strings.Split(infos[1], ",")

				if len(t[tag]) > 0 {
					t[tag] = append(t[tag], urlArr...)
				} else {
					t[tag] = urlArr
				}
			}
		}
	}

	for k, v := range t {
		sort.Strings(v)
		t[k] = uniq(v)
	}
	bArr := []string{}

	for k, v := range t {
		rl := len(v)
		if rl > 10 {
			rl = 10
		}

		tmp := strings.Join(v[:rl], ",")

		bArr = append(bArr, "('"+k+"', '"+tmp+"', '"+strconv.Itoa(m[k])+"', '"+anaDate+"', '"+autils.GetCurrentData(time.Now())+"')")
	}

	sqlStr := "INSERT INTO domain (domain, urls, url_count, ana_date, edit_date) VALUES " + strings.Join(bArr, ",")
	_, err = db.Exec(sqlStr)

	autils.ErrHadle(err)
}

// []string数据去重
func uniq(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
