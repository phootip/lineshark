package controller

import (
	"context"
	"encoding/json"
	"fmt"
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
)

// InitSpreadSheetClient init the sheet
func InitSpreadSheetClient() *sheets.Service {
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
	return Sheet
}

// InitSheetRoute init routing
func InitSheetRoute(g *echo.Group) {
	g.GET("/report", handlerGetReport)
} 

func handlerGetReport(c echo.Context) error {
	spreadSheetID := os.Getenv("SPREEDSHEET_ID")
	readRange := "Sheet1!A1:B"
	resp, err := Sheet.Spreadsheets.Values.Get(spreadSheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	result := ""
	if len(resp.Values) == 0 {
		log.Println("No data found.")
	} else {
		log.Println("Showing read result:")
		for _, row := range resp.Values {
			// Print columns A and E, which correspond to indices 0 and 4.
			result = fmt.Sprintf("%s, %s\n", row[0], row[1])
		}
	}
	return c.String(http.StatusOK, result)
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

func renewToken()
