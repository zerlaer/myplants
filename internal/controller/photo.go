package controller

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"myplants/internal/config"
	"myplants/internal/database"
	"myplants/internal/model"
	"myplants/internal/response"
)

// ListPhotos 照片列表
func ListPhotos(c *gin.Context) {
	plantID := c.Query("plant_id")
	var photos []model.Photo
	query := database.DB.Order("taken_at DESC")
	if plantID != "" {
		query = query.Where("plant_id = ?", plantID)
	}
	if err := query.Find(&photos).Error; err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, photos)
}

// UploadPhoto 上传照片
func UploadPhoto(c *gin.Context) {
	plantID := c.PostForm("plant_id")
	if plantID == "" {
		response.Fail(c, "缺少 plant_id")
		return
	}
	remark := c.PostForm("remark")
	takenAtStr := c.PostForm("taken_at")

	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "请上传文件: "+err.Error())
		return
	}

	cfg := config.Get()
	if file.Size > cfg.Upload.MaxSize {
		response.Fail(c, "文件过大")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		response.Fail(c, "不支持的文件格式")
		return
	}

	pid, _ := strconv.Atoi(plantID)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	relPath := fmt.Sprintf("/uploads/plant_%d/%s", pid, filename)
	if err := saveImageFile(file, relPath); err != nil {
		response.Fail(c, "保存文件失败: "+err.Error())
		return
	}
	takenAt := time.Now()
	if takenAtStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", takenAtStr); err == nil {
			takenAt = t
		}
	}

	photo := model.Photo{
		PlantID: uint(pid),
		Path:    relPath,
		Remark:  remark,
		TakenAt: takenAt,
	}
	if err := database.DB.Create(&photo).Error; err != nil {
		response.Fail(c, "保存记录失败: "+err.Error())
		return
	}

	// 如果植物没有头像,自动设置
	var plant model.Plant
	if err := database.DB.First(&plant, pid).Error; err == nil && plant.Avatar == "" {
		plant.Avatar = relPath
		database.DB.Save(&plant)
	}

	response.OKWithMsg(c, "上传成功", photo)
}

// UploadAvatar 上传植物头像
func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "请上传文件")
		return
	}
	cfg := config.Get()
	if file.Size > cfg.Upload.MaxSize {
		response.Fail(c, "文件过大")
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		response.Fail(c, "不支持的文件格式")
		return
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	relPath := fmt.Sprintf("/uploads/avatars/%s", filename)
	if err := saveImageFile(file, relPath); err != nil {
		response.Fail(c, "保存失败: "+err.Error())
		return
	}
	response.OKWithMsg(c, "上传成功", gin.H{"path": relPath})
}

// DeletePhoto 删除照片
func DeletePhoto(c *gin.Context) {
	id := c.Param("id")
	var photo model.Photo
	if err := database.DB.First(&photo, id).Error; err != nil {
		response.Fail(c, "照片不存在")
		return
	}
	// 删除存储中的文件(含 R2 缩略图)
	deleteImageFile(photo.Path)
	database.DB.Delete(&photo)
	response.OKWithMsg(c, "删除成功", nil)
}

// UpdatePhoto 更新照片备注
func UpdatePhoto(c *gin.Context) {
	id := c.Param("id")
	var photo model.Photo
	if err := database.DB.First(&photo, id).Error; err != nil {
		response.Fail(c, "照片不存在")
		return
	}
	var req struct {
		Remark  string     `json:"remark"`
		TakenAt *time.Time `json:"taken_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	photo.Remark = req.Remark
	if req.TakenAt != nil {
		photo.TakenAt = *req.TakenAt
	}
	database.DB.Save(&photo)
	response.OKWithMsg(c, "更新成功", photo)
}
