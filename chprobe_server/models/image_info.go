package models

type ImageInfo struct {
	ID          int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	HostUUID    string `json:"hostUUID" gorm:"column:host_uuid;index"`
	ImageID     string `json:"imageID" gorm:"column:image_id"`
	RepoTags    string `json:"repoTags" gorm:"column:repo_tags"`
	RepoDigests string `json:"repoDigests" gorm:"column:repo_digests"`
	Size        int64  `json:"size" gorm:"column:size"`
	Created     int64  `json:"created" gorm:"column:created"`
	Os          string `json:"os" gorm:"column:os"`
	Architecture string `json:"architecture" gorm:"column:architecture"`
	DockerVersion string `json:"dockerVersion" gorm:"column:docker_version"`
	LastUpdateTime int64 `json:"lastUpdateTime" gorm:"column:last_update_time"`
}

func (i ImageInfo) TableName() string {
	return "image_info"
}
