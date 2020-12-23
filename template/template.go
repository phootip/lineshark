package template

import (
	"encoding/json"

	"github.com/line/line-bot-sdk-go/linebot"
)

// ConfirmTemplate generator
func ConfirmTemplate(data map[string]string) *linebot.FlexMessage {
	jsonString, _ := json.Marshal(data)
	var contents []linebot.FlexComponent
	var footerContents []linebot.FlexComponent

	contents = append(contents, &linebot.TextComponent{
		Text:   "ยืนยันข้อมูลการโอนเงิน?",
		Size: 	"lg",
		Align: 	"center",
		Weight: "bold",
	})
	contents = append(contents, &linebot.TextComponent{
		Margin: "md",
		Text:   "ที่ดินแปลงที่ " + data["parcel"],
	})
	contents = append(contents, &linebot.TextComponent{
		Text:   data["dateStr"],
	})
	contents = append(contents, &linebot.TextComponent{
		Text:   "จำนวนเงิน " + data["amount"],
	})
	footerContents = append(footerContents, &linebot.ButtonComponent{
		Style: 	"primary",
		Action: linebot.NewPostbackAction("ยืนยัน", string(jsonString), "", "ยืนยัน"),
	})
	
	body := linebot.BoxComponent{
		Layout:   linebot.FlexBoxLayoutTypeVertical,
		Contents: contents,
	}
	footer := linebot.BoxComponent{
		Layout:   linebot.FlexBoxLayoutTypeHorizontal,
		Contents: footerContents,
	}
	bubble := linebot.BubbleContainer{
		Body: &body,
		Footer: &footer,
	}
	return linebot.NewFlexMessage("FlexWithCode", &bubble)// Reply Message
}
