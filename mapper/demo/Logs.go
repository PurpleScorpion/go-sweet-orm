package mapperDemo

type Logs struct {
	Id          int32  `json:"id" tableId:"id"`
	CreatedBy   int32  `json:"createdBy"`
	LogLevel    string `json:"logLevel"`
	LogTime     string `json:"logTime"`
	LogType     string `json:"logType"`
	LogContent  string `json:"logContent"`
	CreatedDate string `json:"createdDate"`
}

func (Logs) TableName() string {
	return "logs"
}
