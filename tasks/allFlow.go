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

type flowStruct struct {
	STime string `json:"sBeginTime"`
	ETime string `json:"sEndTime"`
	SType string `json:"sType"`
	SSite string `json:"sSite"`
}

// 请求流量数据
func UpdateAllFlow(db *sql.DB) {
	yesterday := autils.GetCurrentData(time.Now().AddDate(0, 0, -1))
	yesStr := strings.Replace(yesterday, "-", "", -1)

	someTime := autils.GetCurrentData(time.Now().AddDate(0, 0, -2))
	timeStr := strings.Replace(someTime, "-", "", -1)
	fs := flowStruct{
		timeStr,
		yesStr,
		"day",
		"all",
	}

	queryStr, err := json.Marshal(fs)
	autils.ErrHadle(err)

	url := config.AllFlowUrl + string(queryStr)

	val := flowRequest(url)
	storeData(db, val)
}

// 分析处理数据，保存
func storeData(db *sql.DB, data interface{}) {
	now := autils.GetCurrentData(time.Now())

	flows := data.(*map[string]interface{})
	flowData := (*flows)["data"]

	f := flowData.(map[string]interface{})["data"]
	fl := f.([]interface{})
	var sqlStr []string
	lastFs := fl[len(fl)-1]
	val := lastFs.([]interface{})
	var flowArr []string
	for _, v := range val {
		switch t := v.(type) {
		case string:
			s := v.(string)
			flowArr = append(flowArr, "'"+s+"'")
		case float64:
			f := v.(float64)
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
	sqlStr = append(sqlStr, bf.String())
	sql := "INSERT INTO all_flow (date, click, display, cd_rate, ana_date) VALUES " + strings.Join(sqlStr, ",")
	_, err := db.Exec(sql)
	if err == nil {
		fmt.Println("Update all_flow list successfully.")
	}
}

// 请求处理
func flowRequest(url string) *map[string]interface{} {
	res, err := http.Get(url)
	autils.ErrHadle(err)

	body, err := ioutil.ReadAll(res.Body)
	val := map[string]interface{}{}
	json.Unmarshal(body, &val)
	res.Body.Close()
	autils.ErrHadle(err)

	return &val
}
