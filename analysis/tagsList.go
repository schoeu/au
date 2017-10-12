package analysis

import (
	"os"
	"path/filepath"
	"log"
	"bufio"
	"io"
	"strings"
)

var tagFileName = "./taglist"

func GetTagType(cwd string) map[string]int{
	typeCt := map[string]int{}
	times := 0
	fi, err := os.Open(filepath.Join(cwd, tagFileName))
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
		// content := string(a)
		infos := strings.Split(string(a), ",")
		for _, v := range infos {
			typeCt[v] = times
		}
		times ++
	}

	return typeCt
}
