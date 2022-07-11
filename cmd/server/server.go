package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/MiteshSharma/SlackBot/config"
	"github.com/MiteshSharma/SlackBot/notify"

	"github.com/urfave/negroni"

	"github.com/MiteshSharma/SlackBot/api"
	"github.com/MiteshSharma/SlackBot/logger"
	"github.com/MiteshSharma/SlackBot/metrics"
	"github.com/gorilla/mux"
)

type Server struct {
	Router      *mux.Router
	ServerParam *api.ServerParam
	httpServer  *http.Server
	ServerAPI   *api.ServerAPI
}

func NewServer(appContext context.Context, logger logger.Logger, config *config.Config, notify notify.Notifier) *Server {
	metrics := metrics.NewMetrics()
	router := mux.NewRouter()

	serverParam := api.NewServerParam(logger, metrics, config)

	serverApi := api.NewServerAPI(appContext, router, serverParam, notify)

	server := &Server{
		Router:      router,
		ServerParam: serverParam,
		ServerAPI:   serverApi,
	}

	return server
}

func (s *Server) StartServer() {
	n := negroni.New()

	n.UseHandler(s.Router)

	listenAddr := (":" + s.ServerParam.Config.ServerConfig.Port)
	s.ServerParam.Logger.Debug("Staring server", logger.String("address", listenAddr))
	s.httpServer = &http.Server{
		Handler:      n,
		Addr:         listenAddr,
		ReadTimeout:  s.ServerParam.Config.ServerConfig.ReadTimeout * time.Second,
		WriteTimeout: s.ServerParam.Config.ServerConfig.WriteTimeout * time.Second,
	}

	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil {
			s.ServerParam.Logger.Error("Error starting server ", logger.Error(err))
			return
		}
	}()
}

func (s *Server) StopServer(ctx context.Context) {
	s.httpServer.Shutdown(ctx)

	s.ServerParam.Logger.Info("Stopped server")

	os.Exit(0)
}
