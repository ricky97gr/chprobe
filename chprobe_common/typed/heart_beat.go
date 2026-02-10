package typed

type HeartBeatInfo struct {
	UUID       string `json:"uuid"`
	ReportTime int64  `json:"reportTime"`
}
