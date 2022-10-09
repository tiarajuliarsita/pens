package service

import (
	"github.com/tiarajuliarsita/pens/models"
	"github.com/tiarajuliarsita/pens/repository"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

//create
func (svc *Service) Create(pen models.Pen) error {
	err := svc.repository.CreatePens(pen)
	if err != nil {
		return models.NewInternalServerError(err.Error())
	}
	return nil
}

//update
func (svc *Service) Update(newPen models.Pen) error {
	pen, err := svc.repository.GetPen(newPen.ID)
	if err != nil {
		return models.NewInternalServerError(err.Error())
	}
	err = pen.Exist()
	if err != nil {
		return err

	}
	pen.Name = newPen.Name
	pen.Price = newPen.Price

	err = svc.repository.UpdatePens(pen)
	if err != nil {
		return models.NewInternalServerError(err.Error())
	}
	return nil
}

//get all
func (svc *Service) List() ([]models.Pen, error) {
	pens, err := svc.repository.GetPens()
	if err != nil {
		return nil, err
	}
	return pens, nil
}

//get
func (svc *Service) Get(id int) (models.Pen, error) {
	pen, err := svc.repository.GetPen(id)

	if err != nil {
		return models.Pen{}, err
	}
	if err = pen.Exist(); err != nil {
		return models.Pen{}, err
	}
	return pen, nil
}

//delete
func (svc *Service) Delete(id int) error {

	pen, err := svc.repository.GetPen(id)
	if err != nil {
		return models.NewInternalServerError(err.Error())
	}

	if err = pen.Exist(); err != nil {
		return err

	}

	if err = svc.repository.Delete(pen.ID); err != nil {
		return models.NewInternalServerError(err.Error())
	}
	return nil
}
