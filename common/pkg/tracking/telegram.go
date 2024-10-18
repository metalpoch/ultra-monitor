package tracking

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Telegram struct {
	URL string
}

type payload struct {
	module   string
	category string
	event    string
	message  string
}

func (t Telegram) SendMessage(module, category, event string, err error) {
	dataJson, _ := json.Marshal(payload{module, category, event, err.Error()})
	res, err := http.Post(t.URL, "application/json", bytes.NewBuffer(dataJson))
	if err != nil {
		log.Println("error to send traking on telegram:", err.Error())
		return
	}
	defer res.Body.Close()

	if _, err := io.ReadAll(res.Body); err != nil {
		log.Println("error to read tracking telegram response:", err.Error())
	}
}
