package tracking

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/goccy/go-json"
)

type SmartModule struct {
	URL string
}

type payload struct {
	Module   string `json:"module"`
	Category string `json:"category"`
	Event    string `json:"event"`
	Message  string `json:"message"`
}

func (this SmartModule) SendMessage(module, category, event string, err error) {
	escapedMessage := strings.ReplaceAll(err.Error(), "\"", "\\\"")
	dataJson, err := json.Marshal(payload{Module: module, Category: category, Event: event, Message: escapedMessage})
	if err != nil {
		log.Println("Error al serializar JSON:", err)
		return
	}

	res, err := http.Post(this.URL, "application/json", bytes.NewBuffer(dataJson))
	if err != nil {
		log.Println("Error al enviar tracking a Telegram:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		log.Printf("Error: %s - Response: %s", res.Status, string(body))
		return
	}
}
