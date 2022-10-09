package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tiarajuliarsita/pens/models"
	"github.com/tiarajuliarsita/pens/service"
)

type Controller struct {
	service service.Service
}

func sendresponse(code int, message string, data interface{}, w http.ResponseWriter) {
	resp := models.Response{
		Code:    code,
		Data:    data,
		Message: message,
	}

	write(resp, code, w)

}

func sendresponseError(err error, w http.ResponseWriter) {

	if v, ok := err.(*models.Errors); ok {
		write(v, v.Code, w)
		return
	}
	data := models.NewErrors(http.StatusInternalServerError, "server error", err.Error())
	write(data, http.StatusInternalServerError, w)

}
func write(data interface{}, code int, w http.ResponseWriter) {
	dataByte, _ := json.Marshal(data)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(dataByte)
}
func NewController(service *service.Service) *Controller {
	return &Controller{
		service: *service,
	}
}
func (ctrl *Controller) Create(w http.ResponseWriter, r *http.Request) {
	dataByte, err := io.ReadAll(r.Body)
	if err != nil {
		sendresponseError(models.NewInternalServerError(err.Error()), w)
		return
	}
	defer r.Body.Close()
	var pen models.Pen
	err = json.Unmarshal(dataByte, &pen)
	if err != nil {
		sendresponseError(models.NewInternalServerError(err.Error()), w)
		return
	}

	if err = pen.Validate(); err != nil {
		sendresponseError(err, w)
	}
	err = ctrl.service.Create(pen)
	if err != nil {
		sendresponseError(err, w)
		return
	}
	sendresponse(http.StatusCreated, "success CREATED DATA", nil, w)
	return
}

func (ctrl *Controller) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		sendresponseError(models.NewBadRequest("Error parameter id"),w)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendresponseError(models.NewBadRequest("Cannot convert id param "),w)
		
		return
	}

	dataByte, err := io.ReadAll(r.Body)
	if err != nil {
		sendresponse(http.StatusInternalServerError, "internal server error, fail convert string to int", err.Error(), w)
		return
	}

	defer r.Body.Close()
	var pen models.Pen
	err = json.Unmarshal(dataByte, &pen)
	if err != nil {
		sendresponse(http.StatusInternalServerError, "internal server error", nil, w)
		return
	}

	pen.ID = idInt
	err = ctrl.service.Update(pen)
	if err != nil {
		sendresponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
		return
	}

	sendresponse(http.StatusCreated, "success update", nil, w)
	return
}

func (ctrl *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		sendresponse(http.StatusBadRequest, "bad request, data id params is null", nil, w)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendresponse(http.StatusInternalServerError, "internal server error, fail convert string to int", nil, w)
		return
	}

	err = ctrl.service.Delete(idInt)
	if err != nil {
		sendresponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
		return
	}

	sendresponse(http.StatusOK, "success update", nil, w)
	return
}

func (ctrl *Controller) List(w http.ResponseWriter, r *http.Request) {
	pens, err := ctrl.service.List()
	if err != nil {
		sendresponse(http.StatusBadRequest, "Internal server error, getpens", err.Error(), w)
		return
	}

	sendresponse(http.StatusOK, "succes", pens, w)
	return
}

func (ctrl *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if id == "" {
		sendresponse(http.StatusBadRequest, "bad request, data id params is null", nil, w)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendresponse(http.StatusInternalServerError, "internal server error, fail convert string to int", nil, w)
		return
	}
	pen, err := ctrl.service.Get(idInt)
	if err != nil {
		sendresponse(http.StatusBadRequest, "Internal server error, getpens", err.Error(), w)
		return
	}

	sendresponse(http.StatusOK, "succes", pen, w)

}
