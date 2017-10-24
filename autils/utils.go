package autils

import (
	"../config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"path/filepath"
)

var tempRs = "result"

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
	dirPath := filepath.Join(cwd, tempRs)
	err := os.RemoveAll(dirPath)
	ErrHadle(err)
	mkDirErr := os.MkdirAll(dirPath, 0777)
	ErrHadle(mkDirErr)
	return dirPath
}
