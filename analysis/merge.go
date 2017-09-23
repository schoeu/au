package analysis

import (
	"encoding/json"
	"log"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"fmt"
)

var (
	rs = []rsMapType{}
)

func MergeInfos(msg rsMapType){
	rs = append(rs, msg)
}

func CalcuData(cwd string) {
	std := rs[0]
	tmp := rs[1:]
	for _, v :=  range tmp {
		for i, it := range v {
			if len(std[i]) > 0 {
				std[i] = append(std[i], it...)
			} else {
				std[i] = it
			}
		}
	}

	dir := ensureDir(cwd)

	b, err := json.Marshal(std)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile("_\\d{2}")
	fmt.Println()
	neme := re.ReplaceAllString(fileName, "")
	if e := ioutil.WriteFile(filepath.Join(dir, neme + tempExt), b, 0777); e != nil {
		log.Fatal(e)
	}
}