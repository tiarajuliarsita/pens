package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	// "strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Pen struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//var pens []pen
func getpen(w http.ResponseWriter, r *http.Request) {

}

func sendresponse(code int, message string, data interface{}, w http.ResponseWriter) {
	resp := response{
		Code:    code,
		Data:    data,
		Message: message,
	}
	dataByte, err := json.Marshal(resp)
	if err != nil {
		resp := response{
			Code:    http.StatusInternalServerError,
			Data:    nil,
			Message: "internal server error",
		}
		dataByte, _ = json.Marshal(resp)

	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(dataByte)

}

func remove(slice []Pen, s int) []Pen {
	return append(slice[:1], slice[s+1:]...)
}

var db *sql.DB

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:root@localhost/pens?sslmode=disable")

	if err != nil {
		panic(err.Error())
	}
	if err = db.Ping(); err != nil {
		panic(err.Error())
	}

	if err != nil {
		fmt.Println("err", err)
		panic(err.Error())
	}

	fmt.Println(db == nil)

	defer db.Close()
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/pens", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(123, db == nil)
		rows, err := db.Query("select id, name, price from pens")
		if err != nil {
			sendresponse(http.StatusBadRequest, "internal server eror", nil, w)
		}
		var pens []Pen
		fmt.Println(rows == nil)
		for rows.Next() {
			var pen Pen

			err = rows.Scan(
				&pen.ID,
				&pen.Name,
				&pen.Price,
			)
			if err != nil {
				sendresponse(http.StatusInternalServerError, "internal server", nil, w)
			}
			pens = append(pens, pen)
		}

		sendresponse(http.StatusOK, "succes", pens, w)
		return

	}).Methods(http.MethodGet)

	r.HandleFunc("/api/v1/pens/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		if id == "" {
			sendresponse(http.StatusBadRequest, "bad request, data id params is null", nil, w)
			return
		}

		rows, err := db.Query("select id, name,price from pens where id = $1", id)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error, fail convert string to int", nil, w)
			return
		}

		var pen Pen
		if rows.Next() {
			err = rows.Scan(
				&pen.ID,
				&pen.Name,
				&pen.Price,
			)

			if err != nil {
				sendresponse(http.StatusInternalServerError, "internal server error, fail convert string to int", nil, w)
				return
			}
		}

		if pen.ID == 0 {
			if err != nil {
				sendresponse(http.StatusNotFound, "data not found", nil, w)
				return
			}
		}

		dataByte, err := io.ReadAll(r.Body)
		if err != nil {
			sendresponse(http.StatusBadRequest, "bad request", nil, w)
			return
		}
		defer r.Body.Close()
		var newpen Pen
		err = json.Unmarshal(dataByte, &newpen)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error", nil, w)
			return
		}

		pen.Name = newpen.Name
		pen.Price = newpen.Price

		_, err = db.Exec("UPDATE pens SET name=$2, Price=$3 WHERE id=$1", pen.ID, pen.Name, pen.Price)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error, update pens", err.Error(), w)
			return
		}

		sendresponse(http.StatusCreated, "success update", nil, w)
		return

	}).Methods(http.MethodPut)

	r.HandleFunc("/api/v1/pens", func(w http.ResponseWriter, r *http.Request) {

		dataByte, err := io.ReadAll(r.Body)
		if err != nil {
			sendresponse(http.StatusBadRequest, "bad request", nil, w)
			return
		}
		defer r.Body.Close()
		var pen Pen
		err = json.Unmarshal(dataByte, &pen)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error", nil, w)
			return
		}
		_, err = db.Exec("insert into pens(name,price) values($1,$2)", pen.Name, pen.Price)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error, get pens", nil, w)
			return
		}

		sendresponse(http.StatusCreated, "success", nil, w)
		return

	}).Methods(http.MethodPost)

	r.HandleFunc("/api/v1/pens/{id}", func(w http.ResponseWriter, r *http.Request) {
		
		//id := r.URL.Query().Get("id")
		id := mux.Vars(r)["id"]
		if id == "" {
			sendresponse(http.StatusBadRequest, "bad request, data id params is null", nil, w)
			return
		}

		rows, err := db.Query("select id, name,price from pens where id = $1", id)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error, fail convert string to int", nil, w)
			return
		}

		var pen Pen
		if rows.Next() {
			err = rows.Scan(
				&pen.ID,
				&pen.Name,
				&pen.Price,
			)

			if err != nil {
				sendresponse(http.StatusInternalServerError, "internal server error, fail convert string to int", nil, w)
				return
			}
		}

		if pen.ID == 0 {
			if err != nil {
				sendresponse(http.StatusNotFound, "data not found", nil, w)
				return
			}
		}

		_, err = db.Exec("DELETE FROM pens  WHERE id=$1", pen.ID)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error,delete pens return err", nil, w)
			return
		}
		sendresponse(http.StatusCreated, "success delete", nil, w)
		return
	}).Methods(http.MethodDelete)

	port := "8000"
	fmt.Println("server run on port", port)
	http.ListenAndServe(":"+port, r)
}
