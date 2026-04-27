package service

import (
	"nosocial/internal/dao"
	"nosocial/internal/model"
)

type DrinkService struct {
	drinkDAO    *dao.DrinkDAO
	categoryDAO *dao.DrinkCategoryDAO
}

func NewDrinkService(drinkDAO *dao.DrinkDAO, categoryDAO *dao.DrinkCategoryDAO) *DrinkService {
	return &DrinkService{
		drinkDAO:    drinkDAO,
		categoryDAO: categoryDAO,
	}
}

func (s *DrinkService) GetCategories() ([]*model.DrinkCategory, error) {
	return s.categoryDAO.List()
}

func (s *DrinkService) GetDrinkList(categoryID uint32, status int8, page, size int) ([]*model.Drink, int64, error) {
	offset := (page - 1) * size
	return s.drinkDAO.List(categoryID, status, offset, size)
}

func (s *DrinkService) GetDrinkDetail(id uint64) (*model.Drink, error) {
	return s.drinkDAO.GetByID(id)
}

func (s *DrinkService) GetRecommended() ([]*model.Drink, error) {
	return s.drinkDAO.ListRecommended()
}

// Admin methods
func (s *DrinkService) CreateDrink(drink *model.Drink) error {
	return s.drinkDAO.Create(drink)
}

func (s *DrinkService) UpdateDrink(drink *model.Drink) error {
	return s.drinkDAO.Update(drink)
}

func (s *DrinkService) DeleteDrink(id uint64) error {
	return s.drinkDAO.Delete(id)
}

func (s *DrinkService) CreateCategory(category *model.DrinkCategory) error {
	return s.categoryDAO.Create(category)
}

func (s *DrinkService) UpdateCategory(category *model.DrinkCategory) error {
	return s.categoryDAO.Update(category)
}

func (s *DrinkService) DeleteCategory(id uint32) error {
	return s.categoryDAO.Delete(id)
}
