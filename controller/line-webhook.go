package controller

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// "github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	// Bot linebot singleton
	Bot *linebot.Client
	// TestUser my userid for testing
	TestUser string
	reportTemplate []byte
)

func init() {
	log.Print("Initialize LineBot....")
	var err error
	Bot, err = linebot.New(
		os.Getenv("LINE_SECRET"),
		os.Getenv("LINE_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	TestUser = os.Getenv("TEST_USER")
	reportTemplate, err = ioutil.ReadFile("./template/report.txt")
	if err != nil {
		log.Fatal(err)
	}
}

// HandlerCallback handle line webhook event
func HandlerCallback(c echo.Context) error{
	events, err := Bot.ParseRequest(c.Request())
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.String(400, "error")
		} else {
			c.String(500, "error")
		}
		return err
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			log.Print(event.Source.UserID)
			msg := event.Message.(*linebot.TextMessage)
			log.Print(msg.Text)
		}
	}
	return c.String(http.StatusOK, "received")
}

// LinePushMessage to push message to line
func LinePushMessage(user string, msg string) {
	if _, err := Bot.PushMessage(user, linebot.NewTextMessage(msg)).Do(); err != nil {
		log.Fatal(err)
	}
}

// LinePushReport to push report
func LinePushReport(user string, data map[string]string) error {
	msg := flexMessageFormat(reportTemplate, data)
	container, err := linebot.UnmarshalFlexMessageJSON(msg)
	// err is returned if invalid JSON is given that cannot be unmarshalled
	if err != nil {
		log.Fatal(err)
	}
	if _, err := Bot.PushMessage(
		user,
		linebot.NewFlexMessage("alt text", container),
	).Do(); err != nil {
		return err
	}
	return nil
}

func flexMessageFormat(template []byte, data map[string]string) []byte {
	msg := make([]byte, len(template))
	copy(msg, template)
	for k,v := range data {
		msg = bytes.Replace(msg, []byte("{"+k+"}"),[]byte(v),1)
	}
	return msg
}
