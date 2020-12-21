package controller

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/otiai10/gosseract"
)

// InitAPI add /util to the server
func InitAPI(g *echo.Group) {
	log.Println("Initialize /api ....")
	g.GET("/push", handlerPush)
	g.POST("/sms", handlerSms)
}

// HandlerCallback handle line webhook event
func handlerPush(c echo.Context) error {
	if _, err := Bot.PushMessage(os.Getenv("TEST_USER"), linebot.NewTextMessage("hello phootip from util")).Do(); err != nil {
		log.Fatal(err)
	}
	return c.String(http.StatusOK, "pushed")
}

func handlerSms(c echo.Context) error {
	sender := c.FormValue("sender")
	sms := c.FormValue("sms")
	log.Println("Sender: ", sender)
	log.Println("SMS: ", sms)
	return c.String(http.StatusOK, "Sms received")
}
// มีเงิน 3.00บ.เข้าบ/ชxx5340เหลือ 50,006.00 บ.15/12/20@23:04

// DetectText detect text from image
func detectText(image []byte) string {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImageFromBytes(image)
	client.Languages = []string{"eng","tha"}
	text, err := client.Text()
	if err != nil {
		log.Println(err)
	}
	return text
}

// Temp for testing
func Temp() {
	
}
