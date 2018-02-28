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
		if len(contentSplit) > 1 {
			bf.WriteString("(")
			bf.WriteString("'" + contentSplit[0] + "'," + contentSplit[1])
			bf.WriteString(",")
			bf.WriteString("'" + date + "'")
			bf.WriteString(",")
			bf.WriteString("'" + autils.GetCurrentData(time.Now()) + "'")
			bf.WriteString(")")

			infoArr = append(infoArr, bf.String())
		}
	})

	storeBrowsersData(infoArr, db)
}

func storeBrowsersData(rs []string, db *sql.DB) {
	sqlStr := "INSERT INTO browsers (type, num, date, ana_date) VALUES " + strings.Join(rs, ",")
	_, err := db.Exec(sqlStr)
	autils.ErrHadle(err)
}
