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

type domainStruct struct {
	STime string `json:"sBeginTime"`
	ETime string `json:"sEndTime"`
	SType string `json:"sType"`
	SSite string `json:"sSite"`
}

// 获取单个站点数据
func GetDomains(domain string, db *sql.DB) {
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
		fmt.Println("Update site_flow list successfully.")
	}
}
