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
	yesterday := autils.GetCurrentData(time.Now().AddDate(0, 0, -1))

	flows := data.(*map[string]interface{})
	flowData := (*flows)["data"]

	f := flowData.(map[string]interface{})["data"]
	fl := f.([]interface{})
	var sqlStr []string
	for _, v := range fl {
		val := v.([]interface{})
		var flowArr []string
		var bf bytes.Buffer
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
		flowArr = append(flowArr, "'"+yesterday+"'")
		bf.WriteString("(")
		bf.WriteString(strings.Join(flowArr, ","))
		bf.WriteString(")")
		sqlStr = append(sqlStr, bf.String())
	}
	sql := "INSERT INTO all_flow (date, click, display, cd_rote, ana_date) VALUES " + strings.Join(sqlStr, ",")
	fmt.Println(sql)
	_, err := db.Exec(sql)
	autils.ErrHadle(err)

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