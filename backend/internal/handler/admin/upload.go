package admin

import (
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// allowedDirs 预定义允许的上传目录白名单
var allowedDirs = map[string]bool{
	"banners": true,
	"drinks":  true,
	"tarot":   true,
	"reviews": true,
	"avatars": true,
	"uploads": true,
}

var dirSanitizeRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// allowedExts 可附加在 key 末尾的扩展名白名单（不带点）
var allowedExts = map[string]bool{
	"png": true, "jpg": true, "jpeg": true, "gif": true, "webp": true, "bmp": true,
}

// sanitizeExt 清洗扩展名，非法或不在白名单返回空字符串
func sanitizeExt(ext string) string {
	ext = strings.ToLower(strings.TrimSpace(ext))
	ext = strings.TrimPrefix(ext, ".")
	if ext == "" || !allowedExts[ext] {
		return ""
	}
	return ext
}

type UploadAdminHandler struct {
	uploadService *service.UploadService
}

func NewUploadAdminHandler(uploadService *service.UploadService) *UploadAdminHandler {
	return &UploadAdminHandler{uploadService: uploadService}
}

func (h *UploadAdminHandler) GetOSSSignature(c *gin.Context) {
	dir := sanitizeDir(c.Query("dir"))
	ext := sanitizeExt(c.Query("ext"))

	cred, err := h.uploadService.GetOSSSignature(dir, ext)
	if err != nil {
		response.SafeError(c, err)
		return
	}
	response.Success(c, cred)
}

// FixObjectInline 将已上传对象的 Content-Disposition 设为 inline
func (h *UploadAdminHandler) FixObjectInline(c *gin.Context) {
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

// sanitizeDir 清理并校验上传目录，防止路径遍历攻击
func sanitizeDir(dir string) string {
	dir = strings.TrimSpace(dir)
	dir = strings.Trim(dir, "/")

	// 空值使用默认目录
	if dir == "" {
		return "uploads"
	}

	// 禁止路径遍历字符
	if strings.Contains(dir, "..") || strings.Contains(dir, "\\") || strings.Contains(dir, "/") {
		return "uploads"
	}

	// 只允许字母、数字、下划线和连字符
	if !dirSanitizeRe.MatchString(dir) {
		return "uploads"
	}

	// 检查是否在白名单中
	if !allowedDirs[dir] {
		return "uploads"
	}

	return dir
}
