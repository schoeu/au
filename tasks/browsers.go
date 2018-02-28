package tasks

import (
	"fmt"
	"database/sql"
	"path/filepath"
	"regexp"
	"bytes"
	"strings"
	"time"

	"../autils"
	"../config"
)

func AnaBrowsers(db *sql.DB) {
	cwd := autils.GetCwd()
	fileName := "target_out"
	filePath := filepath.Join(cwd, config.BrowsersPath, fileName)
	fmt.Println(filePath)

	splitReg := regexp.MustCompile("\\t")


	infoArr := []string{}
	
	autils.AnaLogFile(filePath, func (c string) {
		var bf bytes.Buffer

		contentSplit := splitReg.Split(c, -1)

		bf.WriteString("(")
		bf.WriteString(strings.Join(contentSplit, ","))
		bf.WriteString(",")
		bf.WriteString(autils.GetCurrentData(time.Now()))
		bf.WriteString(")")

		infoArr = append(infoArr, bf.String())
	})
	
	fmt.Println(infoArr, db)
}

func storeBrowsersData(rs []string, db *sql.DB) {
	sqlStr := "INSERT INTO browsers (type, num, date) VALUES " + strings.Join(rs, ",")
	_, err := db.Exec(sqlStr)
	autils.ErrHadle(err)
}
