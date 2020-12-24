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
	// log.Println(uploadEvidence("./template/example2.jpg"))
	// data := map[string]string{
	// 	"amount":"10,000.00", 
	// 	"date":"11/29/2020 10:07:00", 
	// 	"dateStr":"29 Nov 2020, 10:07", 
	// 	"parcel":"15",
	// }
	// err := addTransaction(data)
	// log.Println(err)
}
