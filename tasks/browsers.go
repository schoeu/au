package tasks

import (
	"bytes"
	"database/sql"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"../autils"
	"../config"
)

func AnaBrowsers(db *sql.DB, date string) {
	cwd := autils.GetCwd()
	shortDate := strings.Replace(date, "-", "", -1)
	fileName := "target_out"
	filePath := filepath.Join(cwd, config.BrowsersPath, shortDate, fileName)

	splitReg := regexp.MustCompile("\\t")

	infoArr := []string{}

	autils.AnaLogFile(filePath, func(c string) {
		var bf bytes.Buffer

		contentSplit := splitReg.Split(c, -1)

		bf.WriteString("(")
		bf.WriteString(strings.Join(contentSplit, ","))
		bf.WriteString(",")
		bf.WriteString(autils.GetCurrentData(time.Now()))
		bf.WriteString(")")

		infoArr = append(infoArr, bf.String())
	})
}

func storeBrowsersData(rs []string, db *sql.DB) {
	sqlStr := "INSERT INTO browsers (type, num, date) VALUES " + strings.Join(rs, ",")
	_, err := db.Exec(sqlStr)
	autils.ErrHadle(err)
}
