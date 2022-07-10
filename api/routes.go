package api

import (
	"github.com/gorilla/mux"
)

type Router struct {
	Root    *mux.Router // ''
	APIRoot *mux.Router // 'api/v1'
	Setting *mux.Router // 'api/v1/setting'
	Metrics *mux.Router // 'api/v1/metrics'
}

func (a *ServerAPI) setupRoutes() {
	a.Router.Root = a.MainRouter
	a.Router.APIRoot = a.MainRouter.PathPrefix("/api/v1").Subrouter()
	a.Router.Metrics = a.Router.APIRoot.PathPrefix("/metrics").Subrouter()

	a.InitHealthCheck()

	a.Metrics.SetupHttpHandler(a.Router.Metrics)
}
