package application

import (
	"errors"
	"multicliente-backend/internal/features/inventory/article/domain"
	notificationDomain "multicliente-backend/internal/features/notification/domain"
)

type articleService struct {
	repo                domain.ArticleRepository
	notificationService notificationDomain.NotificationService
}

func NewArticleService(repo domain.ArticleRepository, notificationService notificationDomain.NotificationService) domain.ArticleService {
	return &articleService{
		repo:                repo,
		notificationService: notificationService,
	}
}

func (s *articleService) CreateArticle(req *domain.CreateArticleRequest, companyID uint, createdBy *uint) (*domain.Article, error) {
	art := &domain.Article{
		CompanyID:  companyID,
		CategoryID: req.CategoryID,
		Name:       req.Name,
		CreateBy:   createdBy,
	}

	if err := s.repo.Create(art); err != nil {
		return nil, err
	}

	if createdBy != nil && s.notificationService != nil {
		_ = s.notificationService.TriggerEntityEventNotification(companyID, *createdBy, "/inventory/items", "CREATE", art.Name)
	}

	return s.repo.FindByID(art.ID)
}

func (s *articleService) GetArticleByID(id uint) (*domain.Article, error) {
	art, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("article not found")
	}
	return art, nil
}

func (s *articleService) GetArticlesByCompany(companyID uint) ([]domain.Article, error) {
	return s.repo.FindAllByCompany(companyID)
}

func (s *articleService) UpdateArticle(id uint, req *domain.UpdateArticleRequest, updatedBy *uint) (*domain.Article, error) {
	art, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("article not found")
	}

	if req.CategoryID != nil {
		art.CategoryID = *req.CategoryID
	}
	if req.Name != nil {
		art.Name = *req.Name
	}
	art.UpdateBy = updatedBy

	if err := s.repo.Update(art); err != nil {
		return nil, err
	}

	if updatedBy != nil && s.notificationService != nil {
		_ = s.notificationService.TriggerEntityEventNotification(art.CompanyID, *updatedBy, "/inventory/items", "EDIT", art.Name)
	}

	return s.repo.FindByID(art.ID)
}

func (s *articleService) DeleteArticle(id uint) error {
	art, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("article not found")
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	if art.CreateBy != nil && s.notificationService != nil {
		_ = s.notificationService.TriggerEntityEventNotification(art.CompanyID, *art.CreateBy, "/inventory/items", "DELETE", art.Name)
	}
	return nil
}
