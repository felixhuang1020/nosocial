package service

import (
	"nosocial/internal/dao"
	"nosocial/internal/model"
)

type BannerService struct {
	bannerDAO *dao.BannerDAO
}

func NewBannerService(bannerDAO *dao.BannerDAO) *BannerService {
	return &BannerService{bannerDAO: bannerDAO}
}

func (s *BannerService) GetBannersByPosition(position int8) ([]*model.Banner, error) {
	return s.bannerDAO.ListByPosition(position)
}

func (s *BannerService) GetBannerList(page, size int) ([]*model.Banner, int64, error) {
	offset := (page - 1) * size
	return s.bannerDAO.ListAll(offset, size)
}

func (s *BannerService) CreateBanner(banner *model.Banner) error {
	return s.bannerDAO.Create(banner)
}

func (s *BannerService) UpdateBanner(banner *model.Banner) error {
	return s.bannerDAO.Update(banner)
}

func (s *BannerService) DeleteBanner(id uint32) error {
	return s.bannerDAO.Delete(id)
}
