package controller

import (
	"log"
	"net/http"
	"os"

	_ "image/jpeg"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
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
// 30 พ . ย . 2563 - 12:46
//  20,000.00

// DetectText detect text from image

// Temp for testing
func Temp() {
	log.Println("running Temp func....")
	// data := map[string]string{
	// 	"amount":"35,000.00",
	// 	"date":"12/25/2020 17:03:00",
	// 	"dateStr":"25 Dec 2020, 17:03",
	// 	"messageID":"13265349866930",
	// 	"parcel":"15",
	// }
	// addTransaction(data)
}
