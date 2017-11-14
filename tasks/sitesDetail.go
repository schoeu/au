package tasks

import (
	"../autils"
	"../config"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type siteDetail struct {
	Sday string `json:"sDay"`
}

func GetSitesData(db *sql.DB) {
	ss := siteDetail{}
	dayTime := time.Now().AddDate(0, 0, -2)
	start := autils.GetCurrentData(dayTime)
	timeStr := strings.Replace(start, "-", "", -1)

	ss.Sday = timeStr

	queryStr, err := json.Marshal(ss)
	autils.ErrHadle(err)

	rsUrl := config.SiteDetail + string(queryStr)
	rs := requestDetail(rsUrl)

	getSiteDetail(db, rs, start)
}

// 获取站点详细数据
func getSiteDetail(db *sql.DB, detailInfo siteJson, date string) {
	if detailInfo.Retcode != 0 {
		fmt.Println(date, "Update site_detail error.")
		return
	}
	// sj.Data.Data   [][]interface{}
	now := autils.GetCurrentData(time.Now())

	sjData := detailInfo.Data.Data
	var siteInfos []string
	for _, v := range sjData {
		var flowArr []string
		v = v[1:]
		for _, val := range v {
			switch t := val.(type) {
			case string:
				s := val.(string)
				// 兼容接口部分字段
				if s == "-" {
					s = "0"
				}
				// 去掉数据中"%"
				s = strings.Replace(s, "%", "", -1)
				flowArr = append(flowArr, "'"+s+"'")
			case int:
				f := val.(int)
				flowArr = append(flowArr, strconv.Itoa(f))
			case float64:
				f := val.(float64)
				num := strconv.FormatFloat(f, 'f', 2, 64)
				flowArr = append(flowArr, num)
			default:
				_ = t
			}
		}
		flowArr = append(flowArr, "'"+date+"'")
		flowArr = append(flowArr, "'"+now+"'")

		var bf bytes.Buffer
		bf.WriteString("(")
		bf.WriteString(strings.Join(flowArr, ","))
		bf.WriteString(")")
		siteInfos = append(siteInfos, bf.String())
	}

	sql := "INSERT INTO site_detail (domain, total_pv, pv, pv_rate, estimated_pv, estimated_pv_rate, pattern_estimated_pv, urls, record_url, record_rate, pass_url, pass_rate, relative_url, effect_url, effect_pv, ineffect_url, ineffect_pv, shield_url, date, ana_date) VALUES " +
		strings.Join(siteInfos, ",")

	_, err := db.Exec(sql)
	autils.ErrHadle(err)
	if err == nil {
		fmt.Println(date, "Update site_detail successfully.")
	}
}

func requestDetail(url string) siteJson {
	res, err := http.Get(url)
	autils.ErrHadle(err)

	rj := siteJson{}
	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &rj)
	res.Body.Close()
	autils.ErrHadle(err)
	return rj
}
