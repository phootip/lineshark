package controller

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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
// 30 พ . ย . 2563 - 12:46
//  20,000.00

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
	rawData := `@ โอ น เง ิ น ส ํ า เร ็ จ
	30 พ . ย . 2563 - 12:46
	avid WN: 202011300JkFxhHkucr8akDIz
	`
	_ = rawData
	r, _ := regexp.Compile("[0-9]{2,4} ((ม . ค .)|(ก . พ .)|(มี . ค .)|(เม . ย .)|(พ . ค .)|(มิ . ค .)|(ก . ค .)|(ส . ค .)|(ก . ย .)|(ต . ค .)|(พ . ย .)|(ธ . ค .)) [0-9]{2,4}")
	str := r.FindString(rawData)
	log.Println(str)
	rawDate := "30 พ . ย . 2563"
	rawTime := "12:46"
	monthThaiToEng(rawDate, rawTime)
	// writeRange := "A3"
	// var vr sheets.ValueRange
	
	// myval := []interface{}{``}
	// vr.Values = append(vr.Values, myval)
	// _, err := Sheet.Spreadsheets.Values.Update(spreadSheetID, writeRange, &vr).ValueInputOption("RAW").Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve data from sheet. %v", err)
	// 	return err
	// }
	// return c.String(http.StatusOK, "value written")
}

func monthThaiToEng(rawDate string, rawTime string) {
	date := rawDate
	for i := range monthThai {
		date = strings.ReplaceAll(date, monthThai[i], monthEng[i])
	}
	temp := strings.Split(date, " ")
	if len(temp[2]) == 4 {
		temp[2] = temp[2][2:]
	}
	temp2, _ := strconv.Atoi(temp[2])
	temp[2] = strconv.Itoa(temp2 - 43)

	date = strings.Join(temp, " ")
	date += ", " + rawTime
	log.Println(date)
	layout := "2 Jan 06, 15:04"
	t, err := time.Parse(layout, date)
	if err != nil {
			log.Println(err)
	}
	newLayout := "2/1/06 15:04"
	date = t.Format(newLayout)
	log.Println(date)
}
