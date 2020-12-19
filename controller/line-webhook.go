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
	// Bot linebot singleton
	Bot *linebot.Client
	// TestUser my userid for testing
	TestUser string
)

func init() {
	log.Print("Initialize LineBot....")
	var err error
	Bot, err = linebot.New(
		os.Getenv("LINE_SECRET"),
		os.Getenv("LINE_TOKEN"),
	)
	if err != nil {
		panic(err)
	}
	TestUser = os.Getenv("TEST_USER")
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

// LinePushMsg to push message to line
func LinePushMsg(user string, msg string) {
	if _, err := Bot.PushMessage(user, linebot.NewTextMessage(msg)).Do(); err != nil {
		log.Fatal(err)
	}
}

// LinePushReport to push report
func LinePushReport(user string, id int) error{
	container := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeHorizontal,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "Hello,",
				},
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: "World!",
				},
			},
		},
	}
	if _, err := Bot.PushMessage(
		user,
		linebot.NewFlexMessage("alt text", container),
	).Do(); err != nil {
		return err
	}
	return nil
}
