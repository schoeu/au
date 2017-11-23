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

func GetSiteFlow(db *sql.DB, date time.Time) {
	ss := siteStruct{}
	now := time.Now()
	today := autils.GetCurrentData(now)
	yesterday := autils.GetCurrentData(date)
	yesStr := strings.Replace(yesterday, "-", "", -1)

	someTime := autils.GetCurrentData(date.AddDate(0, 0, -1))
	timeStr := strings.Replace(someTime, "-", "", -1)

	ss.STime = timeStr + "00"
	ss.ETime = yesStr + "23"

	queryStr, err := json.Marshal(ss)
	autils.ErrHadle(err)

	rsUrl := config.SitesUrl + string(queryStr)
	rs := getSites(rsUrl)
	sites := rs.Data.Data
	if len(sites) == 0 {
		fmt.Println("No data for domains.")
		return
	}
	var info []string
	for _, v := range sites {
		if strings.Contains(v, ".") {
			getSiteInfo(v, db, date)
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
	if err == nil {
		fmt.Println("Update domains successfully.")
	}
}

// 获取单个站点数据
func getSiteInfo(domain string, db *sql.DB, date time.Time) {
	yesterday := autils.GetCurrentData(date)
	yesStr := strings.Replace(yesterday, "-", "", -1)

	someTime := autils.GetCurrentData(date.AddDate(0, 0, -1))
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
		fmt.Println(yesStr, domain, "update failed.")
		return
	}

	// sj.Data.Data   [][]interface{}
	now := autils.GetCurrentData(time.Now())
	var siteInfos []string
	sjData := sj.Data.Data

	if len(sjData) == 0 {
		fmt.Println("No data for site_flow.")
		return
	}

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
	if err == nil {
		fmt.Println(yesStr, domain, "update successfully.")
	}
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
