package controller

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/measurement/usecase"
	"gorm.io/gorm"
)

type templateController struct {
	Usecase usecase.TemplateUsecase
}

func NewTemplateController(db *gorm.DB, telegram tracking.SmartModule) *templateController {
	return &templateController{Usecase: *usecase.NewTemplateUsecase(db, telegram)}
}

func (t templateController) AddTemplate(template *model.AddTemplate) error {
	return t.Usecase.Add(template)
}

func (t templateController) UpdateTemplate(id uint, template *model.AddTemplate) error {
	return t.Usecase.Update(id, template)
}

func (t templateController) DeleteTemplate(id uint) error {
	return t.DeleteTemplate(id)
}

func (t templateController) ShowAllTemplates(csv bool) error {
	templates, err := t.Usecase.GetAll()
	if err != nil {
		return err
	}

	pretty := table.NewWriter()
	pretty.SetOutputMirror(os.Stdout)
	pretty.AppendHeader(table.Row{
		"ID",
		"Name",
		"Bandwidth OID",
		"Input OID",
		"Output OID",
		"Created at",
		"Updated at",
	})

	for _, template := range templates {
		pretty.AppendRow(table.Row{
			template.ID,
			template.Name,
			template.OidBw,
			template.OidIn,
			template.OidOut,
			template.CreatedAt.Local().Format("2006-01-02 15:04:05"),
			template.UpdatedAt.Local().Format("2006-01-02 15:04:05"),
		})
		pretty.AppendSeparator()
	}

	if csv {
		pretty.RenderCSV()
	} else {
		pretty.Render()
	}

	return nil
}
