package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tiarajuliarsita/pens/controller"
)

func NewRouter(ctrl *controller.Controller) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/pens", ctrl.List).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/pens/{id}", ctrl.Get).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/pens/{id}", ctrl.Update).Methods(http.MethodPut)
	r.HandleFunc("/api/v1/pens", ctrl.Create).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/pens/{id}", ctrl.Delete).Methods(http.MethodDelete)

	return r
}
