package tasks

import (
	"../autils"
	"../config"
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func StepData(db *sql.DB, date string) {
	shortDate := strings.Replace(date, "-", "", -1)
	splitStr := "\001"
	cwd := autils.GetCwd()
	filePath := filepath.Join(cwd, config.StepPath, shortDate)
	fmt.Println(filePath)
	fi, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	bArr := []string{}
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		infoArr := strings.Split(string(a), splitStr)
		bArr = append(bArr, "('"+infoArr[1]+"', '"+strings.Replace(infoArr[2], "//", "", -1)+"', '"+infoArr[3]+"', '"+date+"', '"+autils.GetCurrentData(time.Now())+"')")
	}

	// 数据存储
	storeStepData(bArr)
}

func storeStepData(bArr []string) {
	sqlStr := "INSERT INTO mip_step (type, url, count, date, ana_date) VALUES " + strings.Join(bArr, ",")

	fmt.Println(sqlStr)
	_, err = db.Exec(sqlStr)
	autils.ErrHadle(err)
}
