package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

var (
	// Sheet sigleton
	Sheet *sheets.Service
	spreadSheetID string
)

// InitSpreadSheetClient init the sheet
func init() {
	log.Println("Initialize SpreadSheet....")
	b, err := ioutil.ReadFile("credentials/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	
	client := getClient(config)
	Sheet, err = sheets.New(client)
	spreadSheetID = os.Getenv("SPREEDSHEET_ID")
}

// InitSheetRoute init routing
func InitSheetRoute(g *echo.Group) {
	g.GET("/report/:id", handlerGetReport)
	g.GET("/write", handlerWrite)
} 

func handlerGetReport(c echo.Context) error {
	id := c.Param("id")
	readRange := id+"!K2:O2"
	resp, err := Sheet.Spreadsheets.Values.Get(spreadSheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	data := make(map[string]string)
	data["id"] = id
	if len(resp.Values) == 0 {
		log.Println("No data found.")
	} else {
		data["monthOrder"] = resp.Values[0][0].(string)
		data["month"] = resp.Values[0][1].(string)
		data["expectedAccu"] = resp.Values[0][2].(string)
		data["paidAccu"] = resp.Values[0][3].(string)
		data["overdue"] = resp.Values[0][4].(string)
		if data["overdue"][0] != byte('-') {
			data["overdueColor"] = "#FF0000"
		} else {
			data["overdueColor"] = "#00FF00"
		}
	}
	LinePushReport(TestUser, data)
	return c.String(http.StatusOK, "success")
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
