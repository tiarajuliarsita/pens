package repository

import (
	"database/sql"

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
func (r* Repository)CreatePens(pen models.Pen) (err error) {
	_, err = r.db.Exec("insert into pens(name,price) values($1,$2)", pen.Name, pen.Price)
	if err != nil {
		return
	}
	return nil
}

//get
func (r* Repository) GetPens() (pens []models.Pen, err error) {
	rows, err :=r. db.Query("select id, name, price from pens")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var pen models.Pen

		err = rows.Scan(
			&pen.ID,
			&pen.Name,
			&pen.Price,
		)
		if err != nil {
			return nil, err
		}
		pens = append(pens, pen)
	}
	return pens, nil
}

func  (r* Repository) GetPen(id int) (pen models.Pen, err error) {
	rows, err := r.db.Query("select id, name, price from pens where id=$1", id)
	if err != nil {
		return models.Pen{}, err
	}

	if rows.Next() {

		err = rows.Scan(
			&pen.ID,
			&pen.Name,
			&pen.Price,
		)
		if err != nil {
			return models.Pen{}, err
		}

	}
	return pen, nil
}

//delete
func  (r* Repository) Delete(id int) (err error) {
	_, err = r.db.Exec("DELETE FROM pens  WHERE id=$1", id)
	if err != nil {

		return err
	}
	return nil
}

//update
func  (r* Repository) UpdatePens(pen models.Pen) (err error) {
	_, err = r. db.Exec("update pens set name = $2, price =$3 where id=$1", pen.ID, pen.Name, pen.Price)
	if err != nil {
		return
	}
	return nil
}
