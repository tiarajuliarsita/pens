package repository

import (
	"database/sql"
	"errors"

	"github.com/tiarajuliarsita/pens/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	//validate
	return &Repository{
		db: db,
	}
}

//create
func (r *Repository) CreatePens(pen models.Pen) (err error) {
	_, err = r.db.Exec("insert into pens(name,price) values($1,$2)", pen.Name, pen.Price)
	if err != nil {
		return errors.New("create pens retuen erroeo: " + err.Error())
	}
	return nil
}

//get all
func (r *Repository) GetPens() (pens []models.Pen, err error) {
	rows, err := r.db.Query("select id, name, price from pens")
	if err != nil {
		return nil, errors.New("get pen retuen erroeo: " + err.Error())
	}

	for rows.Next() {
		var pen models.Pen
		err = rows.Scan(
			&pen.ID,
			&pen.Name,
			&pen.Price,
		)
		if err != nil {
			return nil, errors.New("get pen retuen erroeo: " + err.Error())
		}
		pens = append(pens, pen)

	}
	return pens, nil
}

//get
func (r *Repository) GetPen(id int) (pen models.Pen, err error) {
	rows, err := r.db.Query("select id, name, price from pens where id=$1", id)
	if err != nil {
		return models.Pen{}, errors.New("get pen retuen erroeo: " + err.Error())
	}

	if rows.Next() {

		err = rows.Scan(
			&pen.ID,
			&pen.Name,
			&pen.Price,
		)
		if err != nil {
			return models.Pen{}, errors.New("get pen retuen erroeo: " + err.Error())
		}

	}
	return pen, nil
}

//delete
func (r *Repository) Delete(id int) (err error) {
	_, err = r.db.Exec("DELETE FROM pens  WHERE id=$1", id)
	if err != nil {

		return errors.New("delete pens return erroeo: " + err.Error())
	}
	return nil
}

//update
func (r *Repository) UpdatePens(pen models.Pen) (err error) {
	_, err = r.db.Exec("update pens set name = $2, price =$3 where id=$1", pen.ID, pen.Name, pen.Price)
	if err != nil {
		return
	}
	return nil
}
