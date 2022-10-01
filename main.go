package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type pen struct {
	Name  string `json:"name"`
	Price int `json:"price"`
}

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//var pens []pen

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

func remove(slice []pen, s int) []pen {
	return append(slice[:1], slice[s+1:]...)
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:root@localhost/pens?sslmode=disable")

	if err !=nil{
		panic(err.Error())
	}
	if err=db.Ping();err !=nil{
		panic(err.Error())
	}

	if err != nil {
		fmt.Println("err", err)
		panic(err.Error())
	}
	fmt.Println("her")

	defer db.Close()

	http.HandleFunc("/api/v1/pens", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Println(db == nil)
			rows, err := db.Query("select name, price from pens")
			if err != nil {
				sendresponse(http.StatusBadRequest, "internal server eror", nil, w)
			}
			var pens []pen
			fmt.Println(rows == nil)
			for rows.Next() {
				var pen pen

				err = rows.Scan(
					&pen.Name,
					&pen.Price,
				)
				if err != nil {
					sendresponse(http.StatusInternalServerError, "internal server", nil, w)
				}
				pens = append(pens, pen)
			}
			fmt.Println("here")
			sendresponse(http.StatusOK, "succes", pens, w)
			return
		}
		if r.Method == http.MethodPost {
			dataByte, err := io.ReadAll(r.Body)
			if err != nil {
				sendresponse(http.StatusBadRequest, "bad request", nil, w)
				return
			}
			defer r.Body.Close()
			var pen pen
			err = json.Unmarshal(dataByte, &pen)
			if err != nil {
				sendresponse(http.StatusInternalServerError, "internal server error", nil, w)
				return
			}
			_, err = db.Exec("insert into pens(name,price) values($1,$2)",pen.Name, pen.Price)
			if err != nil {
				sendresponse(http.StatusInternalServerError, "internal server error, get pens", nil, w)
				return
			}
			// pens = append(pens, pen)
			sendresponse(http.StatusCreated, "success", nil, w)
			return
		}

		if r.Method == http.MethodPut {
			id := r.URL.Query().Get("id")

			if id == "" {
				sendresponse(http.StatusBadRequest, "bad request, data id params is null", nil, w)
				return
			}
			//cek id ada atau tidak
			idInt, err := strconv.Atoi(id)
			if err != nil {
				sendresponse(http.StatusInternalServerError, "internal server error, fail convert string to int", nil, w)
				return
			}
			// found:= idInt<=len(pens)
			// if!found {
			// 	sendresponse(http.StatusInternalServerError,"internal error", nil, w)
			// 	return
			// }

			idInt -= 1
			dataByte, err := io.ReadAll(r.Body)
			if err != nil {
				sendresponse(http.StatusBadRequest, "bad request", nil, w)
			}
			defer r.Body.Close()
			var pen pen
			err = json.Unmarshal(dataByte, &pen)
			if err != nil {
				sendresponse(http.StatusInternalServerError, "internal server error", nil, w)
			}

			// pens[idInt].Name = pen.Name
			// pens[idInt].Price = pen.Price

			sendresponse(http.StatusCreated, "success update", nil, w)
			return
		}

		if r.Method == http.MethodDelete {
			id := r.URL.Query().Get("id")

			if id == "" {
				sendresponse(http.StatusBadRequest, "bad request, data id params is null", nil, w)
				return
			}
			//cek id ada atau tidak
			idInt, err := strconv.Atoi(id)
			if err != nil {
				sendresponse(http.StatusInternalServerError, "internak server error, fail convert string to int", nil, w)
				return
			}
			// found:= idInt<=len(pens)
			// if!found {
			// 	sendresponse(http.StatusInternalServerError,"internal error", nil, w)
			// 	return
			// }

			idInt -= 1
			// pens = remove(pens, idInt)
			sendresponse(http.StatusCreated, "success delete", nil, w)
			return
		}

	})

	port := "8000"
	fmt.Println("server run on port", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("err", err)
	}
}
