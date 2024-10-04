package tracking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/model"
)

type Telegram struct {
	BotID  string
	ChatID string
}

func (t Telegram) Notification(module, category, event string, err error) {
	url := fmt.Sprintf(constants.TELEGRAM_API_URL, t.BotID)
	text := fmt.Sprintf(constants.TELEGRAM_HTML_MESSAGE, module, category, event, err)
	jsonValue, _ := json.Marshal(model.Telegram{
		ChatID:                t.ChatID,
		ParseMode:             "HTML",
		Text:                  text,
		DisableWebPagePreview: true,
	})

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println("error to send traking on telegram:", err.Error())
		return
	}
	defer res.Body.Close()

	if _, err := io.ReadAll(res.Body); err != nil {
		log.Println("error to read tracking telegram response:", err.Error())
	}
}
