package analysis

import (
	"path/filepath"
	"io/ioutil"
	"strings"
	"os"
	"bufio"
	"io"
	"sort"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
	"fmt"
	"time"
	"strconv"
	"bytes"
)

var (
	tempDir = "./__au_temp__"
	rsPath string
	db *sql.DB
)

type uniqInfoType map[string][]string

func MergeInfos(cwd string, msg rsMapType) {
	var bf bytes.Buffer
	for k, v := range msg {
		l := len(v)
		if l > 10 {
			l = 10
		}
		bf.WriteString(k)
		bf.WriteString(" ")
		bf.WriteString(strings.Join(v[:l], ","))
		bf.WriteString(strings.Join(v[:l], "\n"))
	}

	rsPath = ensureDir(filepath.Join(cwd, tempDir))
	os.RemoveAll(rsPath)
	finalPath := filepath.Join(rsPath, fileName + tempExt)
	fmt.Printf("\nMerge file in %v\n", finalPath)
	if e := ioutil.WriteFile(finalPath, []byte(bf.String()), 0777); e != nil {
		log.Fatal(e)
	}
}

func CalcuUniqInfo(cwd string) {
	t := uniqInfoType{}

	files, err := ioutil.ReadDir(rsPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fi, err := os.Open(filepath.Join(rsPath, file.Name()))
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
			content := string(a)
			infos := strings.Split(content, " ")
			if len(infos) > 1 {
				tag := infos[0]
				urlArr := strings.Split(infos[1], ",")

				if len(t[tag]) > 0 {
					t[tag] = append(t[tag], urlArr...)
				} else {
					t[tag] = urlArr
				}
			}
		}
	}

	for k, v := range t {
		sort.Strings(v)
		t[k] = uniq(v)
	}
	bArr := []string{}
	n := time.Now().String()
	for k, v := range t {
		rl := len(v)
		l := rl
		if rl > 10 {
			rl = 10
		}

		tmp := strings.Join(v[:rl], ",")

		bArr = append(bArr, "('"+k+"', '"+ tmp+"', '"+ strconv.Itoa(l)   +"', '"+n+"', '"+n+"')")
	}

	openDb(cwd)
	sqlStr := "INSERT INTO domain (domain, urls, url_count, ana_date, edit_date) VALUES " + strings.Join(bArr, ",")
	rs, err := db.Exec(sqlStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rs)

	defer db.Close()
}


func uniq(a []string) (ret []string){
	a_len := len(a)
	for i:=0; i < a_len; i++{
		if (i > 0 && a[i-1] == a[i]) || len(a[i])==0{
			continue;
		}
		ret = append(ret, a[i])
	}
	return
}

func openDb(cwd string) {
	config, err := ioutil.ReadFile(filepath.Join(cwd, "config"))
	if err != nil {
		log.Fatal(err)
	}

	mDb, err := sql.Open("mysql",
		string(config))
	db = mDb

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil{
		log.Fatal(err)
	}
}