package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"google.golang.org/api/sheets/v4"
)

var (
	// Sheet sigleton
	Sheet *sheets.Service
	photoAPI *gphotos.Client
	spreadSheetID string
	monthToDiscount map[string]int
	p *message.Printer
	ctx context.Context
)

// InitSpreadSheetClient init the sheet
func init() {
	ctx = context.Background()
	log.Println("Initialize SpreadSheet....")
	b, err := ioutil.ReadFile("credentials/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	
	client := getClient(config)
	Sheet, _ = sheets.New(client)
	photoAPI, _ = gphotos.NewClient(client)
	spreadSheetID = os.Getenv("SPREEDSHEET_ID")
	monthToDiscount = map[string]int{
		"40เดือน": 50000,
		"50เดือน": 25000,
	}
	p = message.NewPrinter(language.English)
}

// InitSheetRoute init routing
func InitSheetRoute(g *echo.Group) {
	g.GET("/write", handlerWrite)
}

// GetReport generate monthly report
func GetReport(id string) ([]byte, error) {
	values := getSheetValues(id)

	data := make(map[string]string)
	if len(values) == 0 {
		log.Println("No data found.")
	} else {
		data["id"] = values[0][0].(string)
		data["monthOrder"] = values[0][2].(string)
		data["month"] = values[0][3].(string)
		data["expectedAccu"] = values[0][4].(string)
		data["paidAccu"] = values[0][5].(string)
		data["overdue"] = values[0][6].(string)
		if data["overdue"][0] != byte('-') && data["overdue"][0] != byte('0') {
			data["overdueColor"] = "#FF0000"
		} else {
			data["overdueColor"] = "#00FF00"
		}
	}
	log.Println(data)
	report := FlexMessageFormat(reportTemplate, data)
	return report, nil
}

func getFutureReport(id string, month string) []byte{
	discount := monthToDiscount[month]
	values := getSheetValues(id)
	data := make(map[string]string)
	
	data["month"] = month
	data["discount"] = p.Sprint(discount)
	data["id"] = values[0][0].(string)
	data["monthOrder"] = values[0][2].(string)
	data["expectedAccu"] = values[0][4].(string)
	data["paidAccu"] = values[0][5].(string)
	data["goal"] = values[0][7].(string)
	
	goal, _ := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(values[0][7].(string), ".00", ""), ",", ""))
	currentMonth, _ := strconv.Atoi(values[0][2].(string))
	goalMonth, _ := strconv.Atoi(month[:2])
	
	data["newGoal"] = p.Sprint(goal - discount)
	data["remainingMonth"] = p.Sprint(goalMonth - currentMonth)
	data["avgInstallment"] = p.Sprint(float64((goal-discount)/(goalMonth-currentMonth)))
	
	return FlexMessageFormat(futureTemplate, data)
}

func handlerWrite(c echo.Context) error {
	writeRange := "A3"
	var vr sheets.ValueRange
	myval := []interface{}{`=IMAGE("http://finviz.com/fut_chart.ashx?t=ES&p&p=m5&s=m",4,100,200)`}
	vr.Values = append(vr.Values, myval)
	_, err := Sheet.Spreadsheets.Values.Update(spreadSheetID, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
		return err
	}
	return c.String(http.StatusOK, "value written")
}

func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "credentials/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		log.Fatal("SpreadSheet Token Error: ", err)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getSheetValues(id string) [][]interface{} {
	readRange := id+"!A71:H71"
	valueRange, err := Sheet.Spreadsheets.Values.Get(spreadSheetID, readRange).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v\n", err)
		return nil
	}
	return valueRange.Values
}

func addTransaction(data map[string]string) string {
	readRange := fmt.Sprintf("'%v'!J2:L1000", data["parcel"])
	valueRange, err := Sheet.Spreadsheets.Values.Get(spreadSheetID, readRange).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v\n", err)
		return "Something went wrong"
	}
	if contains(valueRange.Values, data["date"]) {
		return "Transaction นี้ได้ถูกบันทึกไปแล้ว"
	}

	url, err := uploadEvidence("temp")

	if err != nil {
		log.Println(err)
		return "Something went wrong"
	}

	valueRange.Values = append(valueRange.Values, []interface{}{data["date"],data["amount"],url})
	_, err = Sheet.Spreadsheets.Values.Update(spreadSheetID, readRange, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Printf("Unable to write data to sheet. %v", err)
		return "Something went wrong"
	}
	return "บันทึกเรียบร้อย"
}

func contains(s [][]interface{}, v string) bool {
	for _, a := range s {
		if a[0] == v {
			return true
		}
	}
	return false
}

func uploadEvidence(filePath string) (string, error) {
	a, err := photoAPI.Albums.List().Do()
	log.Println(err)
	for _, album := range a.Albums {
		log.Println(album.Title, album.Id)
	}
	media, err := photoAPI.AddMediaItem(ctx, filePath, os.Getenv("ALBUM_ID"))
	if err != nil {
		return "", err
	}
	return media.ProductUrl, nil
}
