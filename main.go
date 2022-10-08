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
	"github.com/tiarajuliarsita/pens/models"
	"github.com/tiarajuliarsita/pens/repository"
)

//var pens []pen
func getpen(w http.ResponseWriter, r *http.Request) {

}

func sendresponse(code int, message string, data interface{}, w http.ResponseWriter) {
	resp := models.Response{
		Code:    code,
		Data:    data,
		Message: message,
	}
	dataByte, err := json.Marshal(resp)
	if err != nil {
		resp := models.Response{
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
		pens, err := repository.GetPens(db)
		if err != nil {
			sendresponse(http.StatusBadRequest, "Internal server error, getpens", err.Error(), w)
			return
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

		// rows, err := db.Query("select id, name,price from pens where id = $1", id)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error, fail convert string to int", nil, w)
			return
		}

		pen, err :=repository.GetPen(db, id)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error, fail convert string to int", err.Error(), w)
			return
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
		var newpen models.Pen
		err = json.Unmarshal(dataByte, &newpen)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error", nil, w)
			return
		}

		pen.Name = newpen.Name
		pen.Price = newpen.Price

		err=repository.UpdatePens(db, pen)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
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
		var pen models.Pen
		err = json.Unmarshal(dataByte, &pen)
		if err != nil {

			return
		}
		err = repository.CreatePens(db, pen)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error", err, w)
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
		pen, err := repository.GetPen(db, id)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error,get pen return err", err.Error(), w)
			return
		}

		if pen.ID == 0 {
			if err != nil {
				sendresponse(http.StatusNotFound, "data not found", nil, w)
				return
			}
		}

		err = repository.Delete(db, id)
		if err != nil {
			sendresponse(http.StatusInternalServerError, "internal server error,get pen return err", err.Error(), w)
			return
		}
		sendresponse(http.StatusOK, "success delete", nil, w)
	}).Methods(http.MethodDelete)

	port := "8000"
	fmt.Println("server run on port", port)
	http.ListenAndServe(":"+port, r)
}
