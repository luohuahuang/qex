package protocol

type QEXTestCase struct {
	Product    string  `json:"product"`
	SubProduct string  `json:"sub_product"`
	Service    string  `json:"service"`
	API        string  `json:"api"`
	RunId      string  `json:"run_id"`
	Case       string  `json:"case"`
	Branch     string  `json:"branch"`
	Maintainer string  `json:"maintainer"`
	Timestamp  int64   `json:"timestamp"`
	Duration   float64 `json:"duration"`
	Status     int     `json:"status"`
}

type Git struct {
	RunId      string `json:"run_id"`
	Maintainer string `json:"maintainer"`
	Product    string `json:"product"`
	Case       string `json:"case"`
}

const (
	OKRTypeATSignOff  = 0
	OKRTypeATFoundBug = 1
)

type Jira struct {
	RunId   string `json:"run_id"`
	Product string `json:"product"`
	JiraId  string `json:"jira_id"`
	OKRType int    `json:"okr_type"`
}
