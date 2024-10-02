package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/measurement/model"
	"github.com/metalpoch/olt-blueprint/measurement/repository"
	"gorm.io/gorm"
)

type templateUsecase struct {
	repo     repository.TemplateRepository
	telegram tracking.Telegram
}

func NewTemplateUsecase(db *gorm.DB, telegram tracking.Telegram) *templateUsecase {
	return &templateUsecase{repository.NewTemplateRepository(db), telegram}
}

func (use templateUsecase) Add(template *model.AddTemplate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newTemplate := entity.Template{
		Name:      template.Name,
		OidBw:     template.OidBw,
		OidIn:     template.OidIn,
		OidOut:    template.OidOut,
		PrefixBw:  template.PrefixBw,
		PrefixIn:  template.PrefixIn,
		PrefixOut: template.PrefixOut,
	}

	err := use.repo.Add(ctx, newTemplate)
	if err != nil {
		use.telegram.Notification(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(templateUsecase).Add - use.repo.Add(ctx, %v)", newTemplate),
			err,
		)
	}

	return err
}

func (use templateUsecase) Update(id uint, template *model.AddTemplate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e, err := use.repo.Get(ctx, id)
	if err != nil {
		use.telegram.Notification(
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

	err = use.repo.Update(ctx, e)
	if err != nil {
		use.telegram.Notification(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(templateUsecase).Update - use.repo.Update(ctx, %v)", e),
			err,
		)
	}

	return err
}

func (use templateUsecase) Delete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := use.repo.Delete(ctx, id)
	if err != nil {
		use.telegram.Notification(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(templateUsecase).Delete - use.repo.Delete(ctx, %v)", id),
			err,
		)
	}

	return err
}

func (use templateUsecase) GetByID(id uint) (model.Template, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e, err := use.repo.Get(ctx, id)
	if err != nil {
		use.telegram.Notification(
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
		PrefixBw:  e.PrefixBw,
		PrefixIn:  e.PrefixIn,
		PrefixOut: e.PrefixOut,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}, err
}

func (use templateUsecase) GetAll() ([]model.Template, error) {
	templates := []model.Template{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetAll(ctx)

	if err != nil {
		use.telegram.Notification(
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
			PrefixBw:  e.PrefixBw,
			PrefixIn:  e.PrefixIn,
			PrefixOut: e.PrefixOut,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	return templates, err
}
