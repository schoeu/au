package main

import (
	"./analysis"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	tempRs = "result"
	// 需要分析的日志的类型
	consoleTheme = "%c[0;32;40m%s%c[0m\n"
	dateRe       = regexp.MustCompile("mip_processor.log.(\\d{4}-\\d{2}-\\d{2})")
	fileSize     int64
	anaType      int
	anaPath      string
	anaHelper    string
	helpInfo     string
	pattern      string
	// 一个站点最多保存多个少url
	maxLength int
	logFileRe *regexp.Regexp
	anaDate   string
)

// 主函数
func main() {

	flag.IntVar(&anaType, "type", 1,
		`日志分析类型
	1: 生成域名url列表
	2: 统计组件使用次数
	3: 使用组件的url列表`)
	flag.StringVar(&anaPath, "path", "", "需要分析的日志文件夹的绝对路径")
	flag.StringVar(&anaHelper, "help", helpInfo, "help")
	flag.StringVar(&pattern, "pattern", "mip_processor.log.\\d{4}", "需要统计的日志文件名模式，支持正则，默认为全统计")
	flag.IntVar(&maxLength, "maxLength", 10, "制取默认条数")

	// 获取临时路径
	tmpPath := getCwd()

	flag.Parse()

	if anaPath == "" {
		log.Fatal("")
		return
	}

	logFileRe = regexp.MustCompile(pattern)

	if !filepath.IsAbs(anaPath) {
		anaPath = filepath.Join(tmpPath, "..", anaPath)
	}

	// 清除之前临时文件
	cleanTmp(tmpPath)

	start := time.Now()
	// 读取指定目录下文件list
	readDir(anaPath, tmpPath)

	during := time.Since(start)

	fmt.Printf("File size is %v MB, cost %v\n", fileSize/1048576, during)

	if anaDate == "" {
		anaDate = getCurrentData()
	}

	if anaType == 1 {
		analysis.CalcuUniqInfo(tmpPath, anaDate)
	} else if anaType == 2 {
		analysis.GetTagsMap(tmpPath, anaDate)
	} else if anaType == 3 {
		analysis.GetCountData(tmpPath, anaDate)
	}
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
			if anaDate == "" {
				anaDateArr := dateRe.FindAllStringSubmatch(fileName, -1)
				if len(anaDateArr[0]) > 1 {
					anaDate = anaDateArr[0][1]
				}
			}

			fileSize += file.Size()
			fmt.Printf(consoleTheme, 0x1B, "process[ "+file.Name()+" ]done!", 0x1B)
			fullPath := filepath.Join(path, fileName)
			if anaType == 1 {
				analysis.Process(fullPath, cwd, fileName)
			} else if anaType == 2 {
				analysis.TagsUrl(fullPath, cwd, fileName)
			} else if anaType == 3 {
				analysis.CountData(fullPath)
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

func getCurrentData() string {
	t := time.Now().String()
	return strings.Split(t, " ")[0]
}
