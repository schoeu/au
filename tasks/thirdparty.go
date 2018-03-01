package tasks

import (
	"database/sql"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"../autils"
	"../config"
)

func AnaThirdParty(db *sql.DB, date string) {
	cwd := autils.GetCwd()
	shortDate := strings.Replace(date, "-", "", -1)
	fileName := "sanfang_click_" + shortDate
	filePath := filepath.Join(cwd, config.ThirdPartyPath, fileName)

	splitReg := regexp.MustCompile("\\001")
	var contentSplit []string
	autils.AnaLogFile(filePath, func(c string) {
		contentSplit = splitReg.Split(c, -1)
	})

	contentSplit = append(contentSplit, date, autils.GetCurrentData(time.Now()))
	storeThirdData(contentSplit, db)
}

func storeThirdData(rs []string, db *sql.DB) {
	sqlStr := "INSERT INTO thirdparty (total, filter, date, ana_date) VALUES (" + strings.Join(rs, ",") + ")"
	_, err := db.Exec(sqlStr)
	autils.ErrHadle(err)
}
