package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/repository"
)

type templateUsecase struct {
	repo repository.TemplateRepository
}

func NewTemplateUsecase(db *sql.DB) *templateUsecase {
	return &templateUsecase{repository.NewTemplateRepository(db)}
}

func (use templateUsecase) Add(template *model.AddTemplate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := use.repo.Add(ctx, &entity.Template{
		Name:      template.Name,
		OidBw:     template.OidBw,
		OidIn:     template.OidIn,
		OidOut:    template.OidOut,
		PrefixBw:  template.PrefixBw,
		PrefixIn:  template.PrefixIn,
		PrefixOut: template.PrefixOut,
	})

	// Gestionar errores (con Axios por ejemplo)
	// ...

	return err

}

func (use templateUsecase) GetByID(id uint) (model.Template, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e, err := use.repo.GetByID(ctx, id)
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

	// Gestionar errores (con Axios por ejemplo)
	// ...

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
