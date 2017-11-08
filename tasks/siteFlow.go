package tasks

import (
	"../autils"
	"../config"
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type rsJson struct {
	Data struct {
		Data []string
	}
}

type siteJson struct {
	Retcode int
	Data    struct {
		Data [][]interface{}
	}
}

type siteStruct struct {
	STime string `json:"sBeginTime"`
	ETime string `json:"sEndTime"`
}

func GetSiteFlow(db *sql.DB) {
	ss := siteStruct{}
	now := time.Now()
	today := autils.GetCurrentData(now)
	yesterday := autils.GetCurrentData(now.AddDate(0, 0, -1))
	yesStr := strings.Replace(yesterday, "-", "", -1)

	someTime := autils.GetCurrentData(now.AddDate(0, 0, -2))
	timeStr := strings.Replace(someTime, "-", "", -1)

	ss.STime = timeStr + "00"
	ss.ETime = yesStr + "23"

	queryStr, err := json.Marshal(ss)
	autils.ErrHadle(err)

	rsUrl := config.SitesUrl + string(queryStr)
	rs := getSites(rsUrl)
	sites := rs.Data.Data

	var info []string
	for _, v := range sites {
		if strings.Contains(v, ".") {
			getSiteInfo(v, db)
			var bf bytes.Buffer
			bf.WriteString("(")
			bf.WriteString("'" + v + "', '" + today + "'")
			bf.WriteString(")")
			info = append(info, bf.String())
		}
	}
	updateDomains(info, db)
}

func updateDomains(sites []string, db *sql.DB) {
	_, err := db.Exec("delete from domains")
	autils.ErrHadle(err)

	sql := "INSERT INTO domains (domain, ana_date) VALUES " + strings.Join(sites, ",")
	_, err = db.Exec(sql)
	autils.ErrHadle(err)
}

// 获取单个站点数据
func getSiteInfo(domain string, db *sql.DB) {
	yesterday := autils.GetCurrentData(time.Now().AddDate(0, 0, -1))
	yesStr := strings.Replace(yesterday, "-", "", -1)

	someTime := autils.GetCurrentData(time.Now().AddDate(0, 0, -2))
	timeStr := strings.Replace(someTime, "-", "", -1)
	fs := flowStruct{
		timeStr,
		yesStr,
		"day",
		domain,
	}

	queryStr, err := json.Marshal(fs)
	autils.ErrHadle(err)

	url := config.AllFlowUrl + string(queryStr)

	res, err := http.Get(url)
	autils.ErrHadle(err)

	sj := siteJson{}
	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &sj)
	res.Body.Close()
	autils.ErrHadle(err)

	if sj.Retcode != 0 {
		return
	}

	// sj.Data.Data   [][]interface{}
	now := autils.GetCurrentData(time.Now())
	var siteInfos []string
	sjData := sj.Data.Data
	v := sjData[len(sjData)-1]
	var flowArr []string
	flowArr = append(flowArr, "'"+domain+"'")
	for _, val := range v {
		switch t := val.(type) {
		case string:
			s := val.(string)
			flowArr = append(flowArr, "'"+s+"'")
		case float64:
			f := val.(float64)
			num := strconv.FormatFloat(f, 'f', 4, 64)
			flowArr = append(flowArr, num)
		default:
			_ = t
		}
	}
	flowArr = append(flowArr, "'"+now+"'")

	var bf bytes.Buffer
	bf.WriteString("(")
	bf.WriteString(strings.Join(flowArr, ","))
	bf.WriteString(")")

	siteInfos = append(siteInfos, bf.String())
	sql := "INSERT INTO site_flow (domain, date, click, display, total_click, total_display, cd_rate, flow_rate, ana_date) VALUES " + strings.Join(siteInfos, ",")
	_, err = db.Exec(sql)
	autils.ErrHadle(err)
}

func getSites(url string) rsJson {
	res, err := http.Get(url)
	autils.ErrHadle(err)

	rj := rsJson{}
	body, err := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &rj)
	res.Body.Close()
	autils.ErrHadle(err)
	return rj
}
