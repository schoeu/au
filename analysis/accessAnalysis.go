package analysis

import (
	"../autils"
	"database/sql"
	"strings"
	"time"
)

type domainCt []string

func Access(db *sql.DB) {
	allCh := make(chan []string)
	todayCh := make(chan []string)

	var newCt domainCt

	go allDomain(db, allCh)
	go todayDomain(db, todayCh)
	today := <-todayCh
	all := <-allCh

	for _, v := range today {
		has, _ := autils.HasVal(&all, v)
		if !has {
			newCt = append(newCt, v)
		}
	}
	updateNewDomain(newCt, db)
}

// 近期收集到的站点信息
func allDomain(db *sql.DB, ch chan []string) {
	// 查询范围
	dateRange := 50
	now := time.Now()
	tdby := autils.GetCurrentData(now.AddDate(0, 0, -2))
	lastMonthDay := autils.GetCurrentData(now.AddDate(0, 0, -dateRange))

	rows, err := db.Query("select distinct domain from domain where ana_date > ? and ana_date < ?", lastMonthDay, tdby)
	autils.ErrHadle(err)

	allDc := domainCt{}
	var domain string
	for rows.Next() {
		err := rows.Scan(&domain)
		autils.ErrHadle(err)
		allDc = append(allDc, domain)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	ch <- allDc
}

// 当天收集到的站点信息
func todayDomain(db *sql.DB, ch chan []string) {
	yesterday := autils.GetCurrentData(time.Now().AddDate(0, 0, -1))
	rows, err := db.Query("select distinct domain from domain where ana_date = ?", yesterday)
	autils.ErrHadle(err)

	todayDc := domainCt{}
	var domain string
	for rows.Next() {
		err := rows.Scan(&domain)
		autils.ErrHadle(err)
		todayDc = append(todayDc, domain)
	}
	err = rows.Err()
	autils.ErrHadle(err)

	defer rows.Close()

	ch <- todayDc
}

// 更新新增日期到数据库
func updateNewDomain(d []string, db *sql.DB) {
	yesterday := autils.GetCurrentData(time.Now().AddDate(0, 0, -1))

	for i, v := range d {
		rsVal := strings.Replace(v, "'", "\\'", -1)
		d[i] = "'" + rsVal + "'"
	}
	dStr := strings.Join(d, ",")
	_, err := db.Exec("update domain set access_date = ? where ana_date = ? and domain in (?)", yesterday, yesterday, dStr)
	autils.ErrHadle(err)
}
