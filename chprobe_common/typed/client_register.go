package typed

type ClientRegisterInfo struct {
	ClientType    string `json:"clientType"`
	Hostname      string `json:"hostname"`
	IP            string `json:"ip"`
	MachineID     string `json:"machineId"`
	OsType        string `json:"osType"`
	Os            string `json:"os"`
	Arch          string `json:"arch"`
	KernelVersion string `json:"kernelVersion"`
	Version       string `json:"version"`
	Timestamp     int64  `json:"timestamp"`
}
