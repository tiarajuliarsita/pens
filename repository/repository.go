package repository

import (
	"database/sql"

	"github.com/tiarajuliarsita/pens/models"
)

//create
func CreatePens(db *sql.DB, pen models.Pen) (err error) {
	_, err = db.Exec("insert into pens(name,price) values($1,$2)", pen.Name, pen.Price)
	if err != nil {
		return
	}
	return nil
}

//get
func GetPens(db *sql.DB) (pens []models.Pen, err error) {
	rows, err := db.Query("select id, name, price from pens")
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

func GetPen(db *sql.DB, id string) (pen models.Pen, err error) {
	rows, err := db.Query("select id, name, price from pens where id=$1", id)
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
func Delete(db *sql.DB, id string) (err error) {
	_, err = db.Exec("DELETE FROM pens  WHERE id=$1", id)
	if err != nil {

		return err
	}
	return nil
}

//update
func UpdatePens(db *sql.DB, pen models.Pen) (err error) {
	_, err = db.Exec("update pens set name = $2, price =$3 where id=$1", pen.ID, pen.Name, pen.Price)
	if err != nil {
		return
	}
	return nil
}
