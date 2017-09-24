package main

import (
	"./analysis"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"io/ioutil"
)

var (
	tempRs = "result"
	// 需要分析的日志的类型
	logFileRe   = regexp.MustCompile("mip_processor.log.\\d{4}")
	consoleTheme = "%c[0;32;40m%s%c[0m\n"
	aType = 1
	fileSize int64
)

// 主函数
func main() {

	// 获取临时路径
	tmpPath := getCwd()

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Please input a valid log path.")
		return
	}
	inputPath := args[0]

	if !filepath.IsAbs(inputPath) {
		inputPath = filepath.Join(tmpPath, "..", inputPath)
	}

	// 清除之前临时文件
	cleanTmp(tmpPath)

	start := time.Now()
	// 读取指定目录下文件list
	readDir(inputPath, tmpPath)


	during := time.Since(start)

	fmt.Printf(consoleTheme, 0x1B, "processed all success", 0x1B)
	fmt.Printf("File size is %v MB, cost %v", fileSize/1048576, during)

}


// 读取指定目录
func readDir(path string, cwd string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileName := file.Name()
		if logFileRe.MatchString(fileName) {
			fileSize += file.Size()
			fmt.Printf(consoleTheme, 0x1B, "process[ "+file.Name()+" ]done!", 0x1B)
			fullPath := filepath.Join(path, fileName)
			if aType == 1 {
				analysis.CountData(fullPath)
			} else {
				analysis.Process(fullPath, cwd, fileName)
			}
		}
	}
}

// 获取程序cwd
func getCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// 清除临时文件&文件夹
func cleanTmp(dir string) {
	rsPath := filepath.Join(dir, "../", tempRs)
	err := os.RemoveAll(rsPath)
	if err != nil {
		log.Fatal(err)
	}
}
