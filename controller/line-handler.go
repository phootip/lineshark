package controller

import (
	"encoding/json"
	"log"
	"strings"

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
	img, err := getLineImage(message.ID)
	if err != nil {
		log.Println(err)
	}
	date, amount, ok := getDateTime(img)
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
		"messageID": message.ID,
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
	msg := addTransaction(data)
	LineReplyMessage(event.ReplyToken,msg)
}
