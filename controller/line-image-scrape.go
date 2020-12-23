package controller

import (
	"bytes"
	"image"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/otiai10/gosseract"
)

var (
	monthThai = [12]string{"ม . ค .", "ก . พ .", "มี . ค .", "เม . ย .", "พ . ค .", "มิ . ค .", "ก . ค .", "ส . ค .", "ก . ย .", "ต . ค .", "พ . ย .", "ธ . ค ."}
	monthEng = [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
)

func getDateTime(file []byte) (time.Time, string, bool){
	text := detectText(file)
	rawDate, err := qrToDate(file)
	if err != nil {
		if err.Error() == "NotFoundException: startSize = 0" {
			log.Println("No QRcode found, trying OCR...")
			rawDate = detectDate(text)
			if len(rawDate) == 0 {
				log.Println("No Date, skip")
				return time.Now(), "", false
			}
			rawDate = monthThaiToEng(rawDate)
		} else {
			log.Printf("%v, %t",err, err)
		}
	}
	rawTime := detectTime(text)
	rawAmount := detectAmount(text)
	rawDate += ", " + rawTime
	layout := "20060102, 15:04"
	t, _ := time.Parse(layout, rawDate)
	return t, rawAmount, true
	// TODO: move this comment
	// newLayout := "1/2/2006 15:04"
	// date := t.Format(newLayout)
	// return date, rawAmount, true
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
	layout := "2 Jan 06"
	t, _ := time.Parse(layout, date)
	newLayout := "20060102"
	return t.Format(newLayout)
}
