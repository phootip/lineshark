package controller

import (
	"log"
	"net/http"
	"os"

	// "github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	// Bot shutup
	Bot *linebot.Client
)

// InitLineBot shutup
func InitLineBot() {
	log.Print("Initialize LineBot")
	log.Print(os.Getenv("LINE_SECRET"),os.Getenv("LINE_TOKEN"))
	var err error
	Bot, err = linebot.New(
		os.Getenv("LINE_SECRET"),
		os.Getenv("LINE_TOKEN"),
	)
	if err != nil {
		panic("Can't init linebot")
	}
}

// HandlerCallback handle line webhook event
func HandlerCallback(c echo.Context) error{
	// m := echo.Map{}
	// if err := c.Bind(&m); err != nil {
	// 	return err
	// }
	// pretty.Print(m)
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
