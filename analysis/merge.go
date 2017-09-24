package analysis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

var (
	rs = []rsMapType{}
)

func MergeInfos(msg rsMapType) {
	rs = append(rs, msg)
}

func CalcuData(cwd string) {
	std := rs[0]
	tmp := rs[1:]
	for _, v := range tmp {
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
	finalPath := filepath.Join(dir, re.ReplaceAllString(fileName, "")+tempExt)
	fmt.Printf("\nMerge file in %v\n", finalPath)
	if e := ioutil.WriteFile(finalPath, b, 0777); e != nil {
		log.Fatal(e)
	}
}
