package controller

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ricky97gr/chprobe/chprobe_server/database"
	"github.com/ricky97gr/chprobe/chprobe_server/models"
	"github.com/ricky97gr/chprobe/chprobe_server/response"
	"github.com/ricky97gr/chprobe/chprobe_server/utils"
)

// GetAuthInfo 获取授权信息
func GetAuthInfo(c *gin.Context) {
	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 从服务器信息中获取产品序列号
	var serverInfo models.ServerInfo
	if err := db.First(&serverInfo).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to get server info")
		return
	}

	// 查询授权信息
	var licenses []models.License
	if err := db.Order("id desc").Find(&licenses).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to get license info")
		return
	}

	// 构建响应数据
	var authInfoList []map[string]interface{}

	for _, license := range licenses {
		// 验证授权有效性
		status := "valid"
		if license.ExpireTime < time.Now().UnixMilli() {
			status = "invalid"
		}

		// 尝试验证和解密授权字符串
		_, valid, err := utils.VerifyLicenseString(license.Content)
		if err != nil || !valid {
			status = "invalid"
		}

		authInfo := map[string]interface{}{
			"id":         strconv.FormatInt(license.ID, 10),
			"type":       license.Type,
			"importTime": time.UnixMilli(license.ImportTime).Format("2006-01-02 15:04:05"),
			"expireTime": time.UnixMilli(license.ExpireTime).Format("2006-01-02 15:04:05"),
			"status":     status,
		}

		authInfoList = append(authInfoList, authInfo)
	}

	// 使用服务器信息中的序列号作为产品序列号
	productSerial := serverInfo.Serial

	// 如果服务器信息中没有序列号，使用默认值
	if productSerial == "" {
		productSerial = "CHPROBE-2024-0001"
	}

	responseData := map[string]interface{}{
		"productSerial": productSerial,
		"authInfo":      authInfoList,
	}

	response.Success(c, responseData, 0)
}

// UploadLicense 上传授权文件/字符串
func UploadLicense(c *gin.Context) {
	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 检查是文件上传还是字符串上传
	file, _, err := c.Request.FormFile("file")
	if err == nil {
		// 文件上传
		defer file.Close()

		// 读取文件内容
		content, err := io.ReadAll(file)
		if err != nil {
			response.Failed(c, response.ErrStruct, "Failed to read license file")
			return
		}

		// 验证和解密授权文件内容
		licenseData, valid, err := utils.VerifyLicenseString(string(content))
		if err != nil || !valid {
			response.Failed(c, response.ErrStruct, "Invalid license file: "+err.Error())
			return
		}

		// 解析授权数据
		now := time.Now()
		serial := fmt.Sprintf("CHPROBE-%d-%04d", now.Year(), now.YearDay())
		licenseType := "企业版"
		expireTime := now.AddDate(1, 0, 0).UnixMilli()

		// 从授权数据中获取过期时间（如果有）
		if exp, ok := licenseData["expire"].(float64); ok {
			expireTime = int64(exp) * 1000 // 转换为毫秒
		}

		// 从授权数据中获取类型（如果有）
		if typ, ok := licenseData["type"].(string); ok {
			licenseType = typ
		}

		// 创建授权信息
		license := models.License{
			Serial:     serial,
			Type:       licenseType,
			Content:    string(content),
			ImportTime: time.Now().UnixMilli(),
			ExpireTime: expireTime,
			Status:     "valid",
		}

		// 保存到数据库
		if err := db.Create(&license).Error; err != nil {
			response.Failed(c, response.ErrDB, "Failed to save license info")
			return
		}

		response.Success(c, nil, 0)
		return
	}

	// 字符串上传
	licenseStr := c.PostForm("license")
	if licenseStr == "" {
		// 尝试从 JSON 请求体中获取
		var jsonData map[string]string
		if err := c.ShouldBindJSON(&jsonData); err == nil {
			licenseStr = jsonData["license"]
		}

		if licenseStr == "" {
			response.Failed(c, response.ErrStruct, "License content is required")
			return
		}
	}

	// 验证和解密授权字符串
	licenseData, valid, err := utils.VerifyLicenseString(licenseStr)
	if err != nil || !valid {
		response.Failed(c, response.ErrStruct, "Invalid license string: "+err.Error())
		return
	}

	// 解析授权数据
	now := time.Now()
	serial := fmt.Sprintf("CHPROBE-%d-%04d", now.Year(), now.YearDay())
	licenseType := "企业版"
	expireTime := now.AddDate(1, 0, 0).UnixMilli()

	// 从授权数据中获取过期时间（如果有）
	if exp, ok := licenseData["expire"].(float64); ok {
		expireTime = int64(exp) * 1000 // 转换为毫秒
	}

	// 从授权数据中获取类型（如果有）
	if typ, ok := licenseData["type"].(string); ok {
		licenseType = typ
	}

	// 创建授权信息
	license := models.License{
		Serial:     serial,
		Type:       licenseType,
		Content:    licenseStr,
		ImportTime: time.Now().UnixMilli(),
		ExpireTime: expireTime,
		Status:     "valid",
	}

	// 保存到数据库
	if err := db.Create(&license).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to save license info")
		return
	}

	response.Success(c, nil, 0)
}

// DeleteLicense 删除授权
func DeleteLicense(c *gin.Context) {
	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 获取授权ID
	licenseID := c.Param("id")
	if licenseID == "" {
		response.Failed(c, response.ErrStruct, "License ID is required")
		return
	}

	// 删除授权
	if err := db.Where("id = ?", licenseID).Delete(&models.License{}).Error; err != nil {
		response.Failed(c, response.ErrDB, "Failed to delete license")
		return
	}

	response.Success(c, nil, 0)
}

// GetLicenseDetail 获取授权详情
func GetLicenseDetail(c *gin.Context) {
	// 获取数据库连接
	db, err := database.GetMysqlClient()
	if err != nil {
		response.Failed(c, response.ErrDB, "Failed to get database client")
		return
	}

	// 获取授权ID
	licenseID := c.Param("id")
	if licenseID == "" {
		response.Failed(c, response.ErrStruct, "License ID is required")
		return
	}

	// 查询授权信息
	var license models.License
	if err := db.Where("id = ?", licenseID).First(&license).Error; err != nil {
		response.Failed(c, response.ErrDB, "License not found")
		return
	}

	// 验证授权有效性
	status := "valid"
	if license.ExpireTime < time.Now().UnixMilli() {
		status = "invalid"
	}

	// 尝试验证和解密授权字符串
	licenseData, valid, err := utils.VerifyLicenseString(license.Content)
	if err != nil || !valid {
		status = "invalid"
	}

	authInfo := map[string]interface{}{
		"id":         strconv.FormatInt(license.ID, 10),
		"serial":     license.Serial,
		"type":       license.Type,
		"importTime": time.UnixMilli(license.ImportTime).Format("2006-01-02 15:04:05"),
		"expireTime": time.UnixMilli(license.ExpireTime).Format("2006-01-02 15:04:05"),
		"status":     status,
		"content":    license.Content,
	}

	// 如果授权有效，添加解析后的数据
	if valid && licenseData != nil {
		authInfo["licenseData"] = licenseData
	}

	response.Success(c, authInfo, 0)
}
