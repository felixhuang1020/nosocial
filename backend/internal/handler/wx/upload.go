package wx

import (
	"fmt"
	"nosocial/internal/middleware"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// wxAllowedDirs 小程序可用的 OSS 目录白名单
// 比 admin 更严：只开放用户根据业务应该上传的目录
var wxAllowedDirs = map[string]bool{
	"reviews": true,
	"avatars": true,
}

var wxDirRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// wxAllowedExts 小程序可附加在 key 末尾的扩展名白名单
var wxAllowedExts = map[string]bool{
	"png": true, "jpg": true, "jpeg": true, "gif": true, "webp": true, "bmp": true,
}

func sanitizeWXExt(ext string) string {
	ext = strings.ToLower(strings.TrimSpace(ext))
	ext = strings.TrimPrefix(ext, ".")
	if ext == "" || !wxAllowedExts[ext] {
		return ""
	}
	return ext
}

type UploadHandler struct {
	uploadService *service.UploadService
}

func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

// sanitizeWXDir 清洗小程序上传目录，非法/非白名单回退 reviews
func sanitizeWXDir(dir string) string {
	dir = strings.TrimSpace(dir)
	dir = strings.Trim(dir, "/")
	if dir == "" {
		return "reviews"
	}
	if strings.Contains(dir, "..") || strings.Contains(dir, "\\") || strings.Contains(dir, "/") {
		return "reviews"
	}
	if !wxDirRe.MatchString(dir) {
		return "reviews"
	}
	if !wxAllowedDirs[dir] {
		return "reviews"
	}
	return dir
}

// GetOSSSignature 小程序直传签名（共用 admin 签名逻辑，目录白名单更严）
func (h *UploadHandler) GetOSSSignature(c *gin.Context) {
	dir := sanitizeWXDir(c.Query("dir"))
	ext := sanitizeWXExt(c.Query("ext"))
	cred, err := h.uploadService.GetOSSSignature(dir, ext)
	if err != nil {
		response.SafeError(c, err)
		return
	}
	response.Success(c, cred)
}

// FixObjectInline 将已上传对象的 Content-Disposition 设为 inline
func (h *UploadHandler) FixObjectInline(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		response.BadRequest(c, "key 不能为空")
		return
	}
	if err := h.uploadService.FixObjectInline(key); err != nil {
		response.SafeError(c, err)
		return
	}
	response.Success(c, gin.H{"key": key, "status": "inline"})
}

// Upload 接收微信小程序上传的文件（已废弃：推荐使用 GetOSSSignature + OSS 直传）
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请上传文件")
		return
	}

	// 验证文件类型
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	if !allowedExts[ext] {
		response.BadRequest(c, "不支持的文件类型")
		return
	}

	// 验证文件大小 (最大 10MB)
	if file.Size > 10*1024*1024 {
		response.BadRequest(c, "文件大小不能超过10MB")
		return
	}

	// 生成文件名
	userID := middleware.GetWXUserID(c)
	timestamp := time.Now().Unix()
	newFilename := fmt.Sprintf("wx_%d_%d%s", userID, timestamp, ext)

	// 确保上传目录存在
	uploadDir := "./uploads/wx"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.Error(c, 1, "创建上传目录失败")
		return
	}

	// 保存文件
	dst := filepath.Join(uploadDir, newFilename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.Error(c, 1, "保存文件失败")
		return
	}

	// 返回文件访问 URL（实际生产环境应返回 CDN 或 OSS 地址）
	fileURL := fmt.Sprintf("/uploads/wx/%s", newFilename)
	response.Success(c, map[string]interface{}{
		"url":  fileURL,
		"name": newFilename,
	})
}
