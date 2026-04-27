package service

import (
	"nosocial/internal/pkg/oss"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

func (s *UploadService) GetOSSSignature(dir, ext string) (*oss.PostPolicy, error) {
	return oss.GeneratePostSignature(dir, ext)
}

// FixObjectInline 将已上传对象的 Content-Disposition 修复为 inline，解决浏览器直接访问时自动下载的问题。
func (s *UploadService) FixObjectInline(key string) error {
	return oss.SetObjectInline(key)
}
