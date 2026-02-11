package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_server/pagination"
	"github.com/ricky97gr/chprobe/chprobe_server/response"
)

// ImageUI 镜像数据传输对象
type ImageUI struct {
	ID             string `json:"id"`
	HostUUID       string `json:"hostUUID"`
	ImageID        string `json:"imageID"`
	RepoTags       string `json:"repoTags"`
	RepoDigests    string `json:"repoDigests"`
	Size           int64  `json:"size"`
	Created        int64  `json:"created"`
	Os             string `json:"os"`
	Architecture   string `json:"architecture"`
	DockerVersion  string `json:"dockerVersion"`
	LastUpdateTime int64  `json:"lastUpdateTime"`
}

// 将ImageInfo转换为ImageUI
func toImageUI(image models.ImageInfo) ImageUI {
	return ImageUI{
		ID:             image.UUID,
		HostUUID:       image.HostUUID,
		ImageID:        image.ImageID,
		RepoTags:       image.RepoTags,
		RepoDigests:    image.RepoDigests,
		Size:           image.Size,
		Created:        image.Created,
		Os:             image.Os,
		Architecture:   image.Architecture,
		DockerVersion:  image.DockerVersion,
		LastUpdateTime: image.LastUpdateTime,
	}
}

// 将[]ImageInfo转换为[]ImageUI
func toImageUIList(images []models.ImageInfo) []ImageUI {
	imageUIs := make([]ImageUI, len(images))
	for i, image := range images {
		imageUIs[i] = toImageUI(image)
	}
	return imageUIs
}

// 查询镜像列表
func GetImageList(c *gin.Context) {
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

	// 查询镜像总数
	var total int64
	if err := db.Model(&models.ImageInfo{}).Count(&total).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to count images")
		return
	}

	// 分页查询镜像列表
	var images []models.ImageInfo
	result := db.Scopes(pagination.ParseQuery(pageQuery)).Find(&images)
	if result.Error != nil {
		response.Failed(c, response.ErrDB, "Failed to query images")
		return
	}

	// 转换为DTO
	imageUIs := toImageUIList(images)

	// 返回响应
	response.Success(c, imageUIs, total)
}

// 查询单个镜像详情
func GetImageDetail(c *gin.Context) {
	// 获取镜像ID
	id := c.Param("id")
	if id == "" {
		response.Failed(c, response.ErrStruct, "Invalid image ID")
		return
	}

	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 查询镜像详情
	var image models.ImageInfo
	result := db.Where("id = ?", id).First(&image)
	if result.Error != nil {
		response.Failed(c, response.ErrRecordNotFound, "Image not found")
		return
	}

	// 转换为DTO
	imageUI := toImageUI(image)

	// 返回响应
	response.Success(c, imageUI, 1)
}
