package utilities

type RequestData struct {
	Recipients []string `json:"recipients"`
	Body       string   `json:"body"`
	Subject    string   `json:"subject"`
}

const (
	Token = ""
)
