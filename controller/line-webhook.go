package controller

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/otiai10/gosseract"
)

var (
	// Bot linebot singleton
	Bot *linebot.Client
	// TestUser my userid for testing
	TestUser string
	reportTemplate []byte
	futureTemplate []byte
	knownClient map[string]bool
	monthThai = [12]string{"ม . ค .", "ก . พ .", "มี . ค .", "เม . ย .", "พ . ค .", "มิ . ค .", "ก . ค .", "ส . ค .", "ก . ย .", "ต . ค .", "พ . ย .", "ธ . ค ."}
	monthEng = [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
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
	reportTemplate, err = ioutil.ReadFile("./template/report.txt")
	if err != nil {
		log.Fatal(err)
	}
	futureTemplate, err = ioutil.ReadFile("./template/future-report.txt")
	if err != nil {
		log.Fatal(err)
	}
	TestUser = os.Getenv("TEST_USER")
	knownClient = map[string]bool{
		TestUser: true,
		os.Getenv("TEST_GROUP"): true,
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
		// log.Println("Sender user: ", event.Source.UserID)
		// log.Println("Sender group: ", event.Source.GroupID)
		if !(knownClient[event.Source.UserID] || knownClient[event.Source.GroupID]) {
			log.Println("User not allowed")
			return c.String(http.StatusOK, "Your Line user is not allowed") 
		}
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				handlerTextMessage(message, event.ReplyToken)
			case *linebot.ImageMessage:
				handlerImageMessage(message, event.ReplyToken)
			}
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

// LinePushFlex to push report
func LinePushFlex(user string, msg []byte) error {
	container, err := linebot.UnmarshalFlexMessageJSON(msg)
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

// LineReplyMessage to reply with text
func LineReplyMessage(replyToken string, msg string) {
	if _, err := Bot.ReplyMessage(replyToken, linebot.NewTextMessage(msg)).Do(); err != nil {
		log.Fatal(err)
	}
}

// LineReplyFlex to reply with flex message
func LineReplyFlex(replyToken string, msg []byte) error {
	container, err := linebot.UnmarshalFlexMessageJSON(msg)
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err := Bot.ReplyMessage(
		replyToken,
		linebot.NewFlexMessage("alt text", container),
	).Do(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// FlexMessageFormat litery string interpolation
func FlexMessageFormat(template []byte, data map[string]string) []byte {
	msg := make([]byte, len(template))
	copy(msg, template)
	for k,v := range data {
		msg = bytes.ReplaceAll(msg, []byte("{"+k+"}"),[]byte(v))
	}
	return msg
}

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

func handlerImageMessage(message *linebot.ImageMessage, replyToken string) {
	log.Println("MessageID: ", message.ID)
	content, err := Bot.GetMessageContent(message.ID).Do()
	if err != nil {
		log.Println(err)
	}
	defer content.Content.Close()
	buffer, _ := ioutil.ReadAll(content.Content)
	text := detectText(buffer)
	rawDate, err := qrToDate(buffer)
	if err != nil {
		if err.Error() == "NotFoundException: startSize = 0" {
			log.Println("No QRcode found, trying OCR...")
			rawDate = detectDate(text)
			if len(rawDate) == 0 {
				log.Println("No Date, skip")
				return
			}
			rawDate = monthThaiToEng(rawDate)
		} else {
			log.Printf("%v, %t",err, err)
		}
	}
	rawTime := detectTime(text)
	rawAmount := detectAmount(text)
	log.Println(rawDate)
	log.Println(rawTime)
	log.Println(rawAmount)
	rawDate += ", " + rawTime
	layout := "20060102, 15:04"
	t, err := time.Parse(layout, rawDate)
	if err != nil {
			log.Println(err)
	}
	newLayout := "1/2/2006 15:04"
	date := t.Format(newLayout)
	data := date + " - " + rawAmount
	result := fmt.Sprintf("The scraped data is : %v", data)
	LineReplyMessage(replyToken, result)
}

func qrToDate(file []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(file))
	if err != nil {
		return "", err
	}
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}
	code := result.String()
	temp := findTag(code, "00")
	temp = findTag(temp, "02")
	return temp[:8], nil
}

func detectDate(text string) string{
	dateReg, _ := regexp.Compile("[0-9]{2,4} ((ม . ค .)|(ก . พ .)|(มี . ค .)|(เม . ย .)|(พ . ค .)|(มิ . ค .)|(ก . ค .)|(ส . ค .)|(ก . ย .)|(ต . ค .)|(พ . ย .)|(ธ . ค .)) [0-9]{2,4}")
	return dateReg.FindString(text)
}

func detectTime(text string) string{
	timeReg, _ := regexp.Compile("[0-9]{2}:[0-9]{2}")
	return timeReg.FindString(text)
}

func detectAmount(text string) string {
	amountReg, _ := regexp.Compile("[0-9]{1,3},[0-9]{3}.[0-9]{2}")
	return amountReg.FindString(text)
}

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

func findTag(code string, tag string) string{
	pos := 0
	for i := 0; i < len(code); i++ {
		if code[i:i+2] != tag {
			i += 2
			skip, _ := strconv.Atoi(code[i:i+2])
			i += skip + 1
			continue
		}
		pos = i + 4
		break
	}
	return code[pos:]
}

func monthThaiToEng(rawDate string) string {
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
	log.Println(date)
	layout := "2 Jan 06"
	t, _ := time.Parse(layout, date)
	newLayout := "20060102"
	return t.Format(newLayout)
}
