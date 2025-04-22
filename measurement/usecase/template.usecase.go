package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/measurement/repository"
	"gorm.io/gorm"
)

type TemplateUsecase struct {
	repo     repository.TemplateRepository
	telegram tracking.SmartModule
}

func NewTemplateUsecase(db *gorm.DB, telegram tracking.SmartModule) *TemplateUsecase {
	return &TemplateUsecase{repository.NewTemplateRepository(db), telegram}
}

func (uc TemplateUsecase) Add(template *model.AddTemplate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newTemplate := entity.Template{
		Name:   template.Name,
		OidBw:  template.OidBw,
		OidIn:  template.OidIn,
		OidOut: template.OidOut,
	}

	err := uc.repo.Add(ctx, newTemplate)
	if err != nil {
		uc.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(templateUsecase).Add - use.repo.Add(ctx, %v)", newTemplate),
			err,
		)
	}

	return err
}

func (uc TemplateUsecase) Update(id uint, template *model.AddTemplate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e, err := uc.repo.Get(ctx, id)
	if err != nil {
		uc.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(templateUsecase).Update - use.repo.Get(ctx, %d)", id),
			err,
		)
		return err
	}

	if template.Name != "" {
		e.Name = template.Name
	}
	if template.OidIn != "" {
		e.OidIn = template.OidIn
	}
	if template.OidOut != "" {
		e.OidOut = template.OidOut
	}
	if template.OidBw != "" {
		e.OidBw = template.OidBw
	}

	err = uc.repo.Update(ctx, e)
	if err != nil {
		uc.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(templateUsecase).Update - use.repo.Update(ctx, %v)", e),
			err,
		)
	}

	return err
}

func (uc TemplateUsecase) Delete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := uc.repo.Delete(ctx, id)
	if err != nil {
		uc.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(templateUsecase).Delete - use.repo.Delete(ctx, %v)", id),
			err,
		)
	}

	return err
}

func (uc TemplateUsecase) GetByID(id uint) (model.Template, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e, err := uc.repo.Get(ctx, id)
	if err != nil {
		uc.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(templateUsecase).GetByID - use.repo.Get(ctx, %d)", id),
			err,
		)
	}

	return model.Template{
		ID:        e.ID,
		Name:      e.Name,
		OidBw:     e.OidBw,
		OidIn:     e.OidIn,
		OidOut:    e.OidOut,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}, err
}

func (uc TemplateUsecase) GetAll() ([]model.Template, error) {
	templates := []model.Template{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.GetAll(ctx)
	if err != nil {
		uc.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			"(templateUsecase).GetAll - use.repo.GetAll(ctx)",
			err,
		)
	}

	for _, e := range res {
		templates = append(templates, model.Template{
			ID:        e.ID,
			Name:      e.Name,
			OidBw:     e.OidBw,
			OidIn:     e.OidIn,
			OidOut:    e.OidOut,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	return templates, err
}
