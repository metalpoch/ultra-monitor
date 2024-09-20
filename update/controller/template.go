package controller

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/usecase"
)

type templateHandler struct {
	Usecase usecase.TemplateUsecase
}

func newTemplateHandler(db *sql.DB) *templateHandler {
	return &templateHandler{Usecase: *usecase.NewTemplateUsecase(db)}
}

func AddTemplate(db *sql.DB, template *model.AddTemplate) error {
	return newTemplateHandler(db).Usecase.Add(template)
}

func ShowAllTemplates(db *sql.DB, csv bool) error {
	templates, err := newTemplateHandler(db).Usecase.GetAll()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"ID",
		"Name",
		"Bandwidth OID",
		"Input OID",
		"Output OID",
		"Created at",
		"Updated at",
	})

	for _, template := range templates {
		t.AppendRow(table.Row{
			template.ID,
			template.Name,
			fmt.Sprintf("%s (%s)", template.OidBw, template.PrefixBw),
			fmt.Sprintf("%s (%s)", template.OidIn, template.PrefixIn),
			fmt.Sprintf("%s (%s)", template.OidOut, template.PrefixOut),
			template.CreatedAt.Local().Format("2006-01-02 15:04:05"),
			template.UpdatedAt.Local().Format("2006-01-02 15:04:05"),
		})
		t.AppendSeparator()
	}

	if csv {
		t.RenderCSV()
	} else {
		t.Render()
	}

	return nil
}
