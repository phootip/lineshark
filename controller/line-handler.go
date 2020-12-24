package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/kr/pretty"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/phootip/lineshark/template"
)


func handlerTextMessage(message *linebot.TextMessage, replyToken string) {
	tokenized := strings.Split(message.Text, " ")
	if tokenized[0] != "lineshark" { return }
	switch tokenized[1] {
	case "ดูรายงาน":
		report, err := GetReport(tokenized[2])
		if err != nil {
			LineReplyMessage(replyToken, "Something went wrong")
		}
		LineReplyFlex(replyToken, report)
	case "คำนวณ":
		month, id := tokenized[2], tokenized[3]
		report := getFutureReport(id, month)
		LineReplyFlex(replyToken, report)
	default:
		LineReplyMessage(replyToken, "Unknown command")
	}
}

func handlerImageMessage(message *linebot.ImageMessage, replyToken string, userID string) {
	log.Println("MessageID: ", message.ID)
	content, err := Bot.GetMessageContent(message.ID).Do()
	if err != nil {
		log.Println(err)
	}
	defer content.Content.Close()
	buffer, _ := ioutil.ReadAll(content.Content)
	date, amount, ok := getDateTime(buffer)
	if !ok {
		return
	}
	parcel := clientParcel[userID]
	humanLayout := "2 Jan 2006, 15:04"
	dateStr := date.Format(humanLayout)

	sheetLayout := "1/2/2006 15:04:05"
	dateSheet := date.Format(sheetLayout)
	data := map[string]string{
		"parcel": parcel,
		"dateStr": dateStr,
		"amount": amount,
		"date": dateSheet,
	}
	flexMsg := template.ConfirmTemplate(data)
	if _, err := Bot.ReplyMessage(replyToken, flexMsg).Do(); err != nil {
		log.Println(err)
	}
	_, _, _ = parcel,dateStr, amount
	// Bot.ReplyMessage(replyToken, msg).Do()
}

func handlerPostback(event *linebot.Event) {
	log.Println(event.Postback.Data)
	data := make(map[string]string)
	json.Unmarshal([]byte(event.Postback.Data), &data)
	
	pretty.Print(data)
}
