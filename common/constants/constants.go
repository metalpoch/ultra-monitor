package constants

const (
	DATABASE              string = "olt"
	FORMAT_DATE           string = "2006-01-02 15:04:05"
	TELEGRAM_API_URL      string = "https://api.telegram.org/bot%s/sendMessage"
	TELEGRAM_HTML_MESSAGE string = `<b>Tracker Error</b>
	
    <b>🧩 Module:</b> %s
	
    <b>🗃️ Category:</b> %s

    <b>⚠️ Event:</b> %s

    <b>💬 Message:</b> %v`
)

const (
	MODULE_UPDATE  string = "Measurement"
	MODULE_AUTH    string = "Auth"
	MODULE_REPORT  string = "Report"
	MODULE_TRAFFIC string = "Traffic"
)

const (
	CATEGORY_SNMP     string = "SNMP"
	CATEGORY_DATABASE string = "Database"
	CATEGORY_UTILS    string = "Utilities"
	CATEGORY_OSM      string = "OpenStreetMap"
)

const BASE_FILEPATH string = "./report/data/"
