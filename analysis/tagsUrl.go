package analysis

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type tagsUrlType map[string][]string
type rsType struct {
	list  []string
	count int
}

type tagCountInfo struct {

}

var (
	tagsUrlArr   = tagsUrlType{}
	tagsRsUrlArr = tagsUrlType{}
	tagTempDir   = "./__au_tag_temp__"
	tagRsPath    string
	tagRelReg    = regexp.MustCompile(`["|']`)
)

const tagMax = 10

func TagsUrl(filePath string, cwd string, fileName string) {
	tagsUrlArr = tagsUrlType{}
	tagsRsUrlArr = tagsUrlType{}
	fi, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		tmpStr := string(a)
		if !tagRelReg.MatchString(tmpStr) && tagRe.MatchString(tmpStr) {
			getTags(tmpStr)
		}
	}

	var buf bytes.Buffer
	for k, v := range tagsUrlArr {
		buf.WriteString(k)
		buf.WriteString(" ")
		b, uDArr := getDiffUrls(v)
		l := len(v)
		if l > tagMax {
			l = tagMax
		}

		buf.WriteString(strings.Join(b[:l], ","))
		buf.WriteString(" ")

		for key, val := range uDArr {
			buf.WriteString("," + key + "=" + strconv.Itoa(val))
		}
		buf.WriteString("\n")
	}

	tagRsPath = ensureDir(filepath.Join(cwd, tagTempDir))
	finalPath := filepath.Join(tagRsPath, fileName+tempExt)
	fmt.Printf("\nMerge file in %v\n", finalPath)
	if e := ioutil.WriteFile(finalPath, []byte(buf.String()), 0777); e != nil {
		log.Fatal(e)
	}
}

func getDiffUrls(val []string) ([]string, map[string]int) {
	var uniqUrlArr []string
	normalArr := []string{}
	uniqDomainArr := map[string]int{}
	for _, v := range val {
		d := GetDomain(v).host
		if uniqDomainArr[d] == 0 && len(uniqUrlArr) <= tagMax {
			uniqUrlArr = append(uniqUrlArr, v)
		} else {
			normalArr = append(normalArr, v)
		}
		uniqDomainArr[d] += 1
	}
	a := append(uniqUrlArr, normalArr...)
	return a, uniqDomainArr
}

func getTags(c string) {
	tagsInfo := pluginRe.FindAllStringSubmatch(c, -1)
	if len(tagsInfo) > 0 && len(tagsInfo[0]) > 1 {
		url := tagsInfo[0][1]
		tags := tagsInfo[0][2]
		tagsArr := strings.Split(tags, ", ")

		for _, v := range tagsArr {
			tagsUrlArr[v] = append(tagsUrlArr[v], url)
		}
	}
}

func GetTagsMap(cwd string, anaDate string) {
	tagCountCtt := map[string]string{}
	files, err := ioutil.ReadDir(tagRsPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fi, err := os.Open(filepath.Join(tagRsPath, file.Name()))
		if err != nil {
			log.Fatal(err)
			return
		}
		defer fi.Close()
		br := bufio.NewReader(fi)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			// content := string(a)
			infos := bytes.Split(a, []byte(" "))
			if len(infos) > 2 {
				tag := string(infos[0])
				urlArr := strings.Split(string(infos[1]), ",")
				tagC := infos[2][1:]
				tagCountCtt[tag] = string(tagC)
				if len(tagsRsUrlArr[tag]) > 0 {
					tagsRsUrlArr[tag] = append(tagsRsUrlArr[tag], urlArr...)
				} else {
					tagsRsUrlArr[tag] = urlArr
				}
			}
		}
	}

	tagTypeInfo := GetTagType(cwd)

	jst := len(tagTypeInfo)
	fmt.Println(jst)


	for k, v := range tagsRsUrlArr {
		sort.Strings(v)
		tagsRsUrlArr[k] = uniq(v)
	}
	bArr := []string{}
	for k, v := range tagsRsUrlArr {
		rl := len(v)
		if rl > 10 {
			rl = 10
		}
		tmp := strings.Join(v[:rl], ",")

		tagCountStr := tagCountCtt[k]
		tagCountNum := strings.Split(tagCountStr, ",")
		bArr = append(bArr, "('"+k+"', '"+tmp+"', '0', '"+string(tagCountStr)+"','" + strconv.Itoa(tagTypeInfo[k]) +"','"+ strconv.Itoa(len(tagCountNum)) +"','"+anaDate+"', '"+time.Now().String()+"')")
	}
	openDb(cwd)
	sqlStr := "INSERT INTO tags (tag_name, urls, url_count, tag_count, tag_type, domain_count, ana_date, edit_date) VALUES " + strings.Join(bArr, ",")
	rs, err := db.Exec(sqlStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rs)
	defer db.Close()
}
