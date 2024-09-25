package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type interfaceRepository struct {
	db *gorm.DB
}

func NewInterfaceRepository(db *gorm.DB) *interfaceRepository {
	return &interfaceRepository{db}
}

func (repo interfaceRepository) Upsert(ctx context.Context, element *entity.Interface) error {
	// update the elements:	element.IfName,	element.IfDescr, element.IfAlias, element.UpdatedAt
	return repo.db.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(element).Error
}

func (repo interfaceRepository) Get(ctx context.Context, id uint) (*entity.Interface, error) {
	element := new(entity.Interface)
	err := repo.db.WithContext(ctx).First(element, id).Error
	return element, err
}

func (repo interfaceRepository) GetAll(ctx context.Context) ([]*entity.Interface, error) {
	var interfaces []*entity.Interface
	err := repo.db.WithContext(ctx).Find(&interfaces).Error
	return interfaces, err
}

func (repo interfaceRepository) GetAllByDevice(ctx context.Context, id uint) ([]*entity.Interface, error) {
	var interfaces []*entity.Interface
	err := repo.db.WithContext(ctx).Where("device_id = ?", id).Find(&interfaces).Error
	return interfaces, err
}
