package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	menu_  []string
	jokes_ []string
	games_ = [...]string{"打一架", "亲一下"}
)

const (
	c_zhihu_pfx = "http://zhihu.com"
	c_joke_pfx  = "/collection/37895484?page="
)

func parseZhihuPage(str string) int {
	var i, j, jokes_num = 0, 0, 0
	for i = 0; i < len(str); {
		j = strings.Index(str[i:], "data-entry-url=\"")
		if j == -1 {
			break
		}
		i = i + j + 16
		if k := strings.Index(str[i:], "\">"); k != -1 {
			answer := str[i : i+k]
			jokes_ = append(jokes_, answer)
			i += k
			jokes_num++
		}
	}
	return jokes_num
}

func getZhihuPageJokes(page_num int) int {
	url := fmt.Sprintf("%s%s%d", c_zhihu_pfx, c_joke_pfx, page_num)
	res, err := http.Get(url)
	if err != nil {
		log.Printf("http get failed, %s\n", err.Error())
		return 0
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Printf("read http resp body failed, %s\n", err.Error())
		return 0
	}
	return parseZhihuPage(string(body))
}

func RandomJoke() string {
	if jokes_ == nil || len(jokes_) == 0 {
		return "no jokes"
	}
	index := RandomNum(0, len(jokes_))
	return c_zhihu_pfx + jokes_[index]
}

func RandomGame() string {
	index := RandomNum(0, len(games_))
	return games_[index]
}

func ReloadJokes() int {
	jokes_ = make([]string, 0, 200)
	for i := 1; getZhihuPageJokes(i) > 0; i++ {
	}
	return len(jokes_)
}

func RandomNum(start, end int) int {
	k := end - start
	return rand.Intn(k) + start
}

func ReloadMenu() int {
	// menu_ = make([]string, 0, 200)
	return 0
}

func RandomMenu() string {
	index := RandomNum(0, len(menu_))
	return menu_[index]
}

func CheckAdd(uid string, money int, detail string) bool {
	db, err := sql.Open("mysql", "online:online_wx@/online?charset=utf8")
	if err != nil {
		log.Print("mysql open failed, %s\n", err.Error())
		return false
	}
	defer db.Close()

	stmt_ins, err := db.Prepare("INSERT INTO ol_check(uid, money, detail, time) VALUES( ?, ?, ?, ? )") // ? = placeholder
	if err != nil {
		log.Print("stmt prepare failed, %s\n", err.Error())
		return false
	}
	defer stmt_ins.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmt_ins.Exec(uid, money, detail, time.Now().Unix())
	if err != nil {
		log.Print("stmt exec failed, %s\n", err.Error())
		return false
	}
	return true
}

func GetChecks(uid string, start_time int, end_time int) bool {
	db, err := sql.Open("mysql", "online:online_wx@/online?charset=utf8")
	if err != nil {
		log.Print("mysql open failed, %s\n", err.Error())
		return false
	}
	defer db.Close()

	stmt_out, err := db.Prepare(`SELECT detail, money FROM ol_check WHERE uid = ?
	and time >= start_time and time <= end_time`)
	if err != nil {
		log.Print("stmt prepare failed, %s\n", err.Error())
		return false
	}
	defer stmt_out.Close()

	_, err = stmt_out.Exec(uid, start_time, end_time)
	if err != nil {
		log.Print("stmt exec failed, %s\n", err.Error())
		return false
	}
	return true
}
