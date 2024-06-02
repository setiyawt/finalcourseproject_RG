package api

import (
	"finalcourseproject/service"
	"fmt"
	"net/http"
)

type API struct {
	userService              service.UserService
	sessionService           service.SessionService
	predictionService        service.PredictionService
	electricityUsagesService service.ElectricityUsagesService

	mux *http.ServeMux
}

func NewAPI(userService service.UserService, sessionService service.SessionService, predictionService service.PredictionService, electricityUsagesService service.ElectricityUsagesService) API {
	mux := http.NewServeMux()
	api := API{
		userService,
		sessionService,
		predictionService,
		electricityUsagesService,
		mux,
	}

	mux.Handle("/user/register", api.Post(http.HandlerFunc(api.Register)))
	mux.Handle("/user/login", api.Post(http.HandlerFunc(api.Login)))
	mux.Handle("/user/logout", api.Get(api.Auth(http.HandlerFunc(api.Logout))))

	mux.Handle("/electricityusage/get-all", api.Get(api.Auth(http.HandlerFunc(api.FetchAllElectricityUsages))))
	mux.Handle("/electricityusage/get", api.Get(api.Auth(http.HandlerFunc(api.FetchSElectricityUsagesByID))))
	mux.Handle("/electricityusage/add", api.Post(api.Auth(http.HandlerFunc(api.StoreelectricityUsages))))
	mux.Handle("/electricityusage/update", api.Put(api.Auth(http.HandlerFunc(api.UpdateelectricityUsages))))
	mux.Handle("/electricityusage/delete", api.Delete(http.HandlerFunc(api.DeleteelectricityUsages)))

	return api
}

func (api *API) Handler() *http.ServeMux {
	return api.mux
}

func (api *API) Start() {
	fmt.Println("starting web server at http://localhost:8080")
	http.ListenAndServe(":8080", api.Handler())
}
