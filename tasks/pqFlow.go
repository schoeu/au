package tasks

import (
	"../autils"
	"../config"
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

func GetQPSites(db *sql.DB) {
	ss := siteDetail{}
	endTime := time.Now().AddDate(0, 0, -2)
	start := autils.GetCurrentData(endTime)
	ss.Sday = strings.Replace(start, "-", "", -1)

	queryStr, err := json.Marshal(ss)
	autils.ErrHadle(err)

	rsUrl := config.PQSiteDetail + string(queryStr)
	rs := requestDetail(rsUrl)

	getSiteDetail(db, rs, start)
}
