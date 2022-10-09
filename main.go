package main

import (
	"fmt"
	"net/http"

	// "strconv"

	_ "github.com/lib/pq"
	"github.com/tiarajuliarsita/pens/config"
	"github.com/tiarajuliarsita/pens/controller"
	"github.com/tiarajuliarsita/pens/repository"
	"github.com/tiarajuliarsita/pens/router"
	"github.com/tiarajuliarsita/pens/service"
)

func main() {

	con := config.NewConfig()
	db := config.NewDatabase(&con.DB)

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	controller := controller.NewController(service)

	r := router.NewRouter(controller)

	fmt.Printf("server run on port %s:%s", con.API.BaseUrl, con.API.Port)
	http.ListenAndServe(con.API.BaseUrl+":"+con.API.Port, r)
}
