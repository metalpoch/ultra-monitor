package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"gorm.io/gorm"
)

type templateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) *templateRepository {
	return &templateRepository{db}
}

func (repo templateRepository) Add(ctx context.Context, template entity.Template) error {
	return repo.db.WithContext(ctx).Create(&template).Error
}

func (repo templateRepository) Get(ctx context.Context, id uint) (entity.Template, error) {
	var template entity.Template
	err := repo.db.WithContext(ctx).First(&template, id).Error
	return template, err
}

func (repo templateRepository) GetAll(ctx context.Context) ([]entity.Template, error) {
	var templates []entity.Template
	err := repo.db.WithContext(ctx).Find(&templates).Error
	return templates, err
}
