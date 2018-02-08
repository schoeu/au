package tasks

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"../autils"
	"../config"
	"fmt"
	"strings"
)

type arrival struct {
	STATUS []string
}

func GetArrivalData(db *sql.DB, date time.Time) {

	t := autils.GetCurrentData(date.AddDate(0, 0, -1))

	// 2017/12/20171220
	queryStr := strings.Join(strings.Split(t, "-")[:2], "/") + "/" + strings.Replace(t, "-", "", -1)

	res, err := http.Get(config.WebbUrl + queryStr)
	autils.ErrHadle(err)

	a := arrival{}
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	json.Unmarshal(body, &a)
	res.Body.Close()
	autils.ErrHadle(err)
}
