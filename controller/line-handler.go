package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
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
	newLayout := "2 jan 2006 15:04"
	dateStr := date.Format(newLayout)

	leftBtn := linebot.NewMessageAction("Yes", "Yes clicked")
	rightBtn := linebot.NewMessageAction("No", "No clicked")
	temp := fmt.Sprintf("บันทึกวันที่ %v \n จำนวนเงิน %v \n ให้แปลงที่ %v ?",dateStr,amount,clientParcel[userID])
	template := linebot.NewConfirmTemplate(temp, leftBtn, rightBtn)// New TemplateMessage
	msg := linebot.NewTemplateMessage("Confirm Box.", template)// Reply Message
	Bot.ReplyMessage(replyToken, msg).Do()
}
