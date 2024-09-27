package constants

const (
	DATABASE                     string = "olt"
	FORMAT_DATE                  string = "2006-01-02 15:04:05"
	TELEGRAM_API_URL             string = "https://api.telegram.org/bot%s/sendMessage"
	TELEGRAM_MARKDOWN_V2_MESSAGE string = `*OLT BLUEPRINT*
	
    *🧩 Module:* %s
	
    *🗃️ Category:* %s

    *⚠️ Event:* %s

    *💬 Message:* %v`
)
