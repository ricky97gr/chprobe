package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_server/pagination"
	"github.com/ricky97gr/chprobe/chprobe_server/response"
)

// HostUI 主机数据传输对象
type HostUI struct {
	UUID          string   `json:"uuid"`
	HostName      string   `json:"hostName"`
	IP            []string `json:"ip"`
	OsType        string   `json:"osType"`
	Os            string   `json:"os"`
	Arch          string   `json:"arch"`
	KernelVersion string   `json:"kernelVersion"`
	CPU           string   `json:"cpu"`
	Memory        string   `json:"memory"`
	Storage       string   `json:"storage"`
	RegisterTime  int64    `json:"registerTime"`
	LastHeartTime int64    `json:"lastHeartTime"`
}

// 将HostInfo转换为HostUI
func toHostUI(host models.HostInfo) HostUI {
	return HostUI{
		UUID:          host.UUID,
		HostName:      host.HostName,
		IP:            host.IP,
		OsType:        host.OsType,
		Os:            host.Os,
		Arch:          host.Arch,
		KernelVersion: host.KernelVersion,
		CPU:           host.CPU,
		Memory:        host.Memory,
		Storage:       host.Storage,
		RegisterTime:  host.RegisterTime,
		LastHeartTime: host.LastHeartTime,
	}
}

// 将[]HostInfo转换为[]HostUI
func toHostUIList(hosts []models.HostInfo) []HostUI {
	hostUIs := make([]HostUI, len(hosts))
	for i, host := range hosts {
		hostUIs[i] = toHostUI(host)
	}
	return hostUIs
}

// 查询主机列表
func GetHostList(c *gin.Context) {
	// 获取分页查询参数
	pageQuery, err := pagination.GetPageQuery(c)
	if err != nil {
		response.Failed(c, response.ErrStruct, "Invalid page query parameters")
		return
	}

	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 查询主机总数
	var total int64
	if err := db.Model(&models.HostInfo{}).Count(&total).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to count hosts")
		return
	}

	// 分页查询主机列表
	var hosts []models.HostInfo
	result := db.Scopes(pagination.ParseQuery(pageQuery)).Find(&hosts)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "Failed to query hosts")
		return
	}

	// 转换为DTO
	hostUIs := toHostUIList(hosts)

	// 返回响应
	response.Success(c, hostUIs, total)
}

// 查询单个主机详情
func GetHostDetail(c *gin.Context) {
	// 获取主机UUID
	uuid := c.Param("uuid")
	if uuid == "" {
		response.Failed(c, response.ErrStruct, "Invalid host UUID")
		return
	}

	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 查询主机详情
	var host models.HostInfo
	result := db.Where("uuid = ?", uuid).First(&host)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "Host not found")
		return
	}

	// 转换为DTO
	hostUI := toHostUI(host)

	// 返回响应
	response.Success(c, hostUI, 1)
}
