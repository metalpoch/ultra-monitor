package usecase

import "github.com/metalpoch/olt-blueprint/common/model"

type FatUsecase interface {
	Add(fat *model.NewFat) error
	Get(id uint) (*model.Fat, error)
	GetAll() ([]model.FatResponse, error)
	Delete(id uint) error
}

type ReportUsecase interface {
	Add(rp *model.NewReport) (string, error)
	Get(id string) (*model.Report, error)
	GetReports(query *model.FindReports) ([]*model.ReportResponse, error)
	GetCategories() ([]*string, error)
}
