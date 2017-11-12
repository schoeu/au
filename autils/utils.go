package autils

import (
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
func OpenDb(dbTyepe string, dbStr string) *sql.DB {
	if dbTyepe == "" {
		dbTyepe = "mysql"
	}
	db, err := sql.Open(dbTyepe, dbStr)
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
