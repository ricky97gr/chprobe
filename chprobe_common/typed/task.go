package typed

type Task struct {
	TaskID      string `json:"taskId"`
	TaskName    string `json:"taskName"`
	Timestamp   int64  `json:"timestamp"`
	PluginID    string `json:"pluginID"`
	Md5         string `json:"md5"`
	DownloadURL string `json:"downloadURL"`
}
