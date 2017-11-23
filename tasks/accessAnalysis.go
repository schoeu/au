package tasks

import (
	"../autils"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type domainCt []string

func Access(db *sql.DB, date time.Time) {
	currentDay := autils.GetCurrentData(date.AddDate(0, 0, -1))
	day := autils.GetCurrentData(date.AddDate(0, 0, -2))
	lastS := autils.GetCurrentData(date.AddDate(0, -3, 0))

	rows, err := db.Query("select domain from site_detail where date = '" + currentDay + "' except select distinct  domain from site_detail where date <= '" + day + "' and date >= '" + lastS + "'")
	autils.ErrHadle(err)

	newer := domainCt{}
	var domain string
	for rows.Next() {
		err := rows.Scan(&domain)
		autils.ErrHadle(err)
		newer = append(newer, domain)
	}

	updateNewDomain(newer, db, currentDay)
}

// 更新新增日期到数据库
func updateNewDomain(d []string, db *sql.DB, date string) {
	for i, v := range d {
		rsVal := strings.Replace(v, "'", "\\'", -1)
		d[i] = "'" + rsVal + "'"
	}
	dStr := strings.Join(d, ",")
	_, err := db.Exec("update site_detail set access_date = '" + date + "' where date = '" + date + "' and domain in (" + dStr + ")")
	autils.ErrHadle(err)
	fmt.Println("update site_detail access date successfully.")
}
