package tasks

import (
	"../autils"
	"../config"
	"bufio"
	"database/sql"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func CustomData(db *sql.DB, date string) {
	shortDate := strings.Replace(date, "-", "", -1)
	splitReg := regexp.MustCompile("\\s+")
	cwd := autils.GetCwd()
	filePath := filepath.Join(cwd, config.CustomPath, shortDate)
	fi, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	var bArr, rs []string
	pos := []int{}
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		line := string(a)
		if strings.Contains(line, "TARGETS") {
			bArr = splitReg.Split(line, -1)
			pos = arrIndexOf(bArr, []string{config.CTotal, config.NormalC, config.CustomC})
		}
		if strings.Contains(line, "All") {
			numArr := splitReg.Split(line, -1)
			for _, v := range pos {
				if len(numArr) >= v {
					rs = append(rs, numArr[v])
				}
			}
		}
	}
	// 数据存储
	if len(rs) > 0 {
		storeCustomData(rs, date, db)
	}
}

func arrIndexOf(data []string, sub []string) []int {
	var pos []int
	for i, v := range data {
		for _, val := range sub {
			if v == val {
				pos = append(pos, i)
				break
			}
		}
	}
	return pos
}

func storeCustomData(rs []string, date string, db *sql.DB) {
	sqlStr := "INSERT INTO custom (total, normal, cust, date, ana_date) VALUES (" + strings.Join(rs, ",") + ", '" + date + "','" + autils.GetCurrentData(time.Now()) + "')"
	_, err := db.Exec(sqlStr)
	autils.ErrHadle(err)
}
