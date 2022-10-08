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

func NewController(service *service.Service) *Controller {
	return &Controller{
		service: *service,
	}
}
func (ctrl *Controller) Create(w http.ResponseWriter, r *http.Request ) {
	dataByte, err := io.ReadAll(r.Body)
	if err != nil {
		sendresponse(http.StatusBadRequest, "bad request", nil, w)
		return
	}
	defer r.Body.Close()
	var pen models.Pen
	err = json.Unmarshal(dataByte, &pen)
	if err != nil {
		sendresponse(http.StatusBadRequest, "bad request", nil, w)

	}
	err = ctrl.service.Create(pen)
	if err != nil {
		sendresponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
	}

	sendresponse(http.StatusCreated, "success", nil, w)
	return
}

func (ctrl *Controller) Update(w http.ResponseWriter, r *http.Request ) {
	id := mux.Vars(r)["id"]

	if id == "" {
		sendresponse(http.StatusBadRequest, "bad request, data id params is null", nil, w)
		return
	}

	idInt, err :=strconv.Atoi(id)
	if err != nil {
		sendresponse(http.StatusInternalServerError, "internal server error, fail convert string to int", nil, w)
		return
	}

	dataByte, err :=io.ReadAll(r.Body)
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

	pen.ID=idInt
	err = ctrl.service.Update( pen)
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

	idInt, err :=strconv.Atoi(id)
	if err != nil {
		sendresponse(http.StatusInternalServerError, "internal server error, fail convert string to int", nil, w)
		return
	}


	err = ctrl.service.Delete( idInt)
	if err != nil {
		sendresponse(http.StatusInternalServerError, "internal server error", err.Error(), w)
		return
	}

	sendresponse(http.StatusOK, "success update", nil, w)
	return
}

func (ctrl *Controller) List(w http.ResponseWriter, r *http.Request ) {
	pens, err := ctrl.service.List()
		if err != nil {
			sendresponse(http.StatusBadRequest, "Internal server error, getpens", err.Error(), w)
			return
		}

		sendresponse(http.StatusOK, "succes", pens, w)
		return
}

func (ctrl *Controller) Get(w http.ResponseWriter, r *http.Request ) {
	id := mux.Vars(r)["id"]

	if id == "" {
		sendresponse(http.StatusBadRequest, "bad request, data id params is null", nil, w)
		return
	}

	idInt, err :=strconv.Atoi(id)
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



