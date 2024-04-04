package entity

type contextKey string

const (
	ConstAcceptLanguage   contextKey = "accept-language"
	ConstLayoutDateFormat string     = "2006-01-02"
	Blue                  string     = "#147BD1"
	Green                 string     = "#2BA031"
	Yellow                string     = "#FFAA00"
	Red                   string     = "#DB313D"
	Pink                  string     = "#f72585"
	GrassGreen            string     = "#b5e48c"
	Purple                string     = "#7209b7"
	Brown                 string     = "#a68a64"
	Indigo                string     = "#480ca8"
	Black                 string     = "#000000"
	Orange                string     = "#FF9D00"
	SystemName            string     = "system"

	RegexNumber string = "[0-9]+"
)

type Ping struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}
