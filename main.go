package main

import (
	"./analysis"
	"./autils"
	"./config"
	"./tasks"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"time"
	"strings"
)

var (
	// 需要分析的日志的类型
	consoleTheme = "%c[0;32;40m%s%c[0m\n"
	dateRe       = regexp.MustCompile("mip_processor.log.(\\d{4}-\\d{2}-\\d{2})")
	fileSize     int64
	anaType      int
	anaPath      string
	anaHelper    string
	helpInfo     string
	pattern      string
	logFileRe    *regexp.Regexp
	anaDate      string
	shortForm    = "2006-01-02"
)

// 主函数
func main() {
	var date, isDebug string
	flag.IntVar(&anaType, "type", 1,
		`日志分析类型
	1: 生成域名url列表
	2: 统计组件使用次数
	3: 使用组件的url列表`)
	flag.StringVar(&anaPath, "path", "", "需要分析的日志文件夹的绝对路径")
	flag.StringVar(&anaHelper, "help", helpInfo, "help")
	flag.StringVar(&pattern, "pattern", "mip_processor.log.\\d{4}", "需要统计的日志文件名模式，支持正则，默认为全统计")
	flag.StringVar(&date, "date", "默认日期为两天前（yyyy-MM-dd）", "拉取数据的日期")
	flag.StringVar(&isDebug, "debug", "", "是否开启调试模式")

	flag.Parse()

	rsPQUrl := config.PQFlowUrl
	// 是否开启调试模式
	if isDebug != "" {
		rsPQUrl = config.PQTestUrl
	}

	pqDB := autils.OpenDb("postgres", rsPQUrl)
	defer pqDB.Close()

	if anaType == 4 {
		// 执行任务
		rsDate := time.Now().AddDate(0, 0, -1)
		trimDate := strings.TrimSpace(date)
		dateReg := regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}")
		if dateReg.MatchString(trimDate) {
			rsDate, _ = time.Parse(shortForm, date)
		}
		runTask(pqDB, rsDate)
		return
	} else if anaType == 5 {
		// mip step part.
		tasks.StepData(pqDB, date)
		return
	}

	if anaPath == "" {
		log.Fatal("Invild log path string.")
		return
	}

	logFileRe = regexp.MustCompile(pattern)

	// 获取临时路径
	tmpPath := autils.GetCwd()

	if !filepath.IsAbs(anaPath) {
		anaPath = filepath.Join(tmpPath, "..", anaPath)
	}

	// 读取指定目录下文件list
	readDir(anaPath, tmpPath)

	during := time.Since(time.Now())

	fmt.Printf("File size is %v MB, cost %v\n", fileSize/1048576, during)

	if anaDate == "" {
		anaDate = autils.GetCurrentData(time.Now())
	}

	if anaType == 1 {
		analysis.CalcuUniqInfo(anaDate, pqDB)
	} else if anaType == 2 {
		analysis.GetTagsMap(anaDate, pqDB)
	} else if anaType == 3 {
		analysis.GetCountData(anaDate, pqDB)
	}
}

// 读取指定目录
func readDir(path string, cwd string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	autils.CleanTmp(filepath.Join(cwd, config.TagTempDir))
	autils.CleanTmp(filepath.Join(cwd, config.TempDir))

	for _, file := range files {
		fileName := file.Name()
		if logFileRe.MatchString(fileName) {
			if anaDate == "" {
				anaDateArr := dateRe.FindAllStringSubmatch(fileName, -1)
				if len(anaDateArr[0]) > 1 {
					anaDate = anaDateArr[0][1]
				}
			}
			fmt.Printf(consoleTheme, 0x1B, "process [ "+file.Name()+" ] done!", 0x1B)
			fileSize += file.Size()
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

// PSCE任务列表部分
func runTask(db *sql.DB, date time.Time) {
	// 更新组件列表
	tasks.UpdateTags(db)
	// 更新流量数据
	tasks.UpdateAllFlow(db, date)
	// 单站点数据
	tasks.GetSiteFlow(db, date)
	// 全量站点详情
	tasks.GetQPSites(db, date)
	// 更新最新接入站点信息
	tasks.Access(db, date)
	// 更新webb日志数据
	//tasks.GetArrivalData(db, date)
}
