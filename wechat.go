package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	//"time"
)

const (
	c_token     = "fishioon"
	c_admin_uid = "olM6ms0pCy7zQPSZbkmhalGYEe3o"
	c_msg_tpl   = `<xml>
                    <ToUserName><![CDATA[%s]]></ToUserName>
                    <FromUserName><![CDATA[%s]]></FromUserName>
                    <CreateTime>%d</CreateTime>
                    <MsgType><![CDATA[%s]]></MsgType>
                    <Content><![CDATA[%s]]></Content>
                    <FuncFlag>0</FuncFlag>
                    </xml>`
)

type Wechat struct {
	ToUserName   string `xml:ToUserName`
	FromUserName string `xml:FromUserName`
	MsgType      string `xml:MsgType`
	CreateTime   int64  `xml:CreateTime`
	Content      string `xml:Content`
}

func Sign(token, timestamp, nonce string) string {
	strs := sort.StringSlice{token, timestamp, nonce}
	strs.Sort()

	buf := make([]byte, 0, len(token)+len(timestamp)+len(nonce))
	buf = append(buf, strs[0]...)
	buf = append(buf, strs[1]...)
	buf = append(buf, strs[2]...)

	hashsum := sha1.Sum(buf)
	return hex.EncodeToString(hashsum[:])
}

func wePack(w *Wechat, content string) string {
	res := fmt.Sprintf(c_msg_tpl, w.FromUserName, w.ToUserName, w.CreateTime, w.MsgType, content)
	return res
}

func weUnpack(xmlstr string, w *Wechat) bool {
	if err := xml.Unmarshal([]byte(xmlstr), &w); err != nil {
		log.Printf("unpack wechat message failed, err:%s\n", err.Error())
		return false
	}
	return true
}

func handleCheck(w *Wechat, inputs []string) string {
	if len(inputs) != 3 {
		return "你输入的姿势好像不对，查询账单的正确姿势: 账单 开始日期 结束日期"
	}
	start_arr := strings.Split(inputs[1], ".")
	// end_arr := strings.Split(inputs[2], ".")
	if len(start_arr) != 3 {
		return "你输入的开始日期好像不对，正确姿势: 2016.8.7"
	}
	//start_date := time.Date(strconv.Atoi(start_arr[0]), strconv.Atoi(start_arr[1], start_arr[2], 0, 0, 0, 0, time.Local)
	//end_date := time.Date(end_arr[0], end_arr[1], end_arr[2], 0, 0, 0, 0, time.Local)
	//res := GetChecks(w.FromUserName, start_time, end_time)
	return ""
}

func addCheck(w *Wechat, inputs []string) string {
	if len(inputs) != 3 {
		return "你输入的姿势好像不对，记录账单的正确姿势: zd 金额 备注"
	}
	money, err := strconv.ParseFloat(inputs[1], 64)
	if err != nil {
		return "你输入的金额好像不对，它应该是个数字"
	}
	// 数据库记录的金额单位为分，用整数保存
	res := CheckAdd(w.FromUserName, int(money*100), inputs[2])
	if res == false {
		return "系统出现了一点小故障，快去联系小金鱼"
	}
	return "添加账单成功"
}

func handleRandom(w *Wechat, content string) string {
	result := ""
	if content == "上上签" {
		result = RandomGame()
	} else if content == "菜单" {
		result = RandomMenu()
	} else {
		num := 100
		if s, err := strconv.Atoi(content); err == nil {
			if s > 1 {
				num = s
			}
		}
		result = strconv.Itoa(RandomNum(0, num))
	}
	return result
}

func weHandleMsg(w *Wechat) string {

	result := ""
	text := strings.TrimSpace(w.Content)
	inputs := strings.Split(text, " ")

	switch strings.ToLower(inputs[0]) {
	case "sj":
		if len(inputs) >= 2 {
			result = handleRandom(w, inputs[1])
		}
	case "xh":
		result = RandomJoke()
	case "zd":
		result = addCheck(w, inputs)
	case "账单":
		result = handleCheck(w, inputs)
	case "jj":
		if w.FromUserName == c_admin_uid {
			go func() {
				jokes_num := ReloadJokes()
				log.Printf("reload jokes, jokes num:%d\n", jokes_num)
			}()
			result = "reloading"
		}
	case "菜单":
		result = "menu"
	default:
		log.Printf("invalid command, msg:%s\n", w.Content)
	}
	if result == "" {
		result = "使用说明: 输入'xh'随机出现一个笑话"
	}
	return result
}

func WeCheckSign(signature, timestamp, nonce, echostr string) string {
	res := Sign(c_token, timestamp, nonce)
	if res == signature {
		return echostr
	} else {
		return "sign failed"
	}
}

func WeProcessMsg(xmlstr string) string {
	w := Wechat{}
	res := ""
	if weUnpack(xmlstr, &w) {
		res = weHandleMsg(&w)
	} else {
		res = "invaild wechat message"
	}
	return wePack(&w, res)
}
