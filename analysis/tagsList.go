package analysis

import (
	"strings"
)

func GetTagType(cwd string) map[string]int {
	typeCt := map[string]int{}
	tagListArr := strings.Split(TagsList, "\n")
	for k, v := range tagListArr {
		infos := strings.Split(string(v), ",")
		for _, v := range infos {
			typeCt[v] = k + 1
		}
	}
	return typeCt
}
