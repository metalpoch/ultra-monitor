package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/gorm"
)

type TemplateRepository interface {
	Add(ctx context.Context, template entity.Template) error
	Get(ctx context.Context, id uint) (*entity.Template, error)
	GetAll(ctx context.Context) ([]entity.Template, error)
	Update(ctx context.Context, template *entity.Template) error
	Delete(ctx context.Context, id uint) error
}

type templateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *templateRepository {
	return &templateRepository{db}
}

func (repo templateRepository) Add(ctx context.Context, template entity.Template) error {
	return repo.db.WithContext(ctx).Create(&template).Error
}

func (repo templateRepository) Update(ctx context.Context, template *entity.Template) error {
	return repo.db.WithContext(ctx).Save(template).Error
}

func (repo templateRepository) Delete(ctx context.Context, id uint) error {
	return repo.db.WithContext(ctx).Delete(&entity.Template{ID: id}).Error
}

func (repo templateRepository) Get(ctx context.Context, id uint) (*entity.Template, error) {
	template := new(entity.Template)
	err := repo.db.WithContext(ctx).First(template, id).Error
	return template, err
}

func (repo templateRepository) GetAll(ctx context.Context) ([]entity.Template, error) {
	var templates []entity.Template
	err := repo.db.WithContext(ctx).Find(&templates).Error
	return templates, err
}
