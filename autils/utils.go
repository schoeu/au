package autils

import (
	"../config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
	"time"
)

// 统一错误处理
func ErrHadle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// 创建数据库链接
func OpenDb(cwd string) *sql.DB {
	db, err := sql.Open("mysql", config.DbConfig)
	ErrHadle(err)

	err = db.Ping()
	ErrHadle(err)
	return db
}

// 创建临时文件夹存放中间文件
func EnsureDir(cwd string) string {
	// err := os.RemoveAll(dirPath)
	// ErrHadle(err)
	mkDirErr := os.MkdirAll(cwd, 0777)
	ErrHadle(mkDirErr)
	return cwd
}

// 清除临时文件&文件夹
func CleanTmp(p string) {
	if p == "" {
		return
	}
	err := os.RemoveAll(p)
	if err != nil {
		log.Fatal(err)
	}
}

// 获取程序cwd
func GetCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// 获取当前日期字符串
//func GetCurrentData() string {
//	t := time.Now().String()
//	return strings.Split(t, " ")[0]
//}

// []string indexOf
func HasVal(a *[]string, it string) (bool, string) {
	for _, v := range *a {
		if v == it {
			return true, v
		}
	}
	return false, ""
}

// 获取当前时间字符串
func GetCurrentData(date time.Time) string {
	t := date.String()
	return strings.Split(t, " ")[0]
}
