package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func path(endpoint string) string {
	return fmt.Sprintf("/api/v1%s", endpoint)
}

func NewRESTApiV1(productionMode bool, logger *zap.Logger) *RESTApiV1 {
	restAPI := &RESTApiV1{
		router:         mux.NewRouter(),
		logger:         logger,
		loggerNoStack:  logger.WithOptions(zap.AddStacktrace(zap.DPanicLevel)),
		productionMode: productionMode,
	}

	restAPI.router.HandleFunc("/", restAPI.IndexPage).Methods(http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodHead)

	restAPI.router.HandleFunc(path("/packetloss/start"), restAPI.PacketlossStart).Methods(http.MethodPost)
	restAPI.router.HandleFunc(path("/packetloss/status"), restAPI.PacketlossStatus).Methods(http.MethodGet)
	restAPI.router.HandleFunc(path("/packetloss/stop"), restAPI.PacketlossStop).Methods(http.MethodPost)

	restAPI.router.HandleFunc(path("/bandwidth/start"), restAPI.BandwidthStart).Methods(http.MethodPost)
	restAPI.router.HandleFunc(path("/bandwidth/status"), restAPI.BandwidthStatus).Methods(http.MethodGet)
	restAPI.router.HandleFunc(path("/bandwidth/stop"), restAPI.BandwidthStop).Methods(http.MethodPost)

	restAPI.router.HandleFunc(path("/latency/start"), restAPI.LatencyStart).Methods(http.MethodPost)
	restAPI.router.HandleFunc(path("/latency/status"), restAPI.LatencyStatus).Methods(http.MethodGet)
	restAPI.router.HandleFunc(path("/latency/stop"), restAPI.LatencyStop).Methods(http.MethodPost)

	restAPI.router.HandleFunc(path("/services/status"), restAPI.NetServicesStatus).Methods(http.MethodGet)

	return restAPI
}

func (a *RESTApiV1) Serve(addr, originAllowed string) error {
	http.Handle("/", a.router)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token"})
	originsOk := handlers.AllowedOrigins([]string{originAllowed})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	a.logger.Info(fmt.Sprintf("serving on %s", addr))

	a.server = &http.Server{
		Addr:    addr,
		Handler: handlers.CORS(originsOk, headersOk, methodsOk)(a.router),
	}

	return a.server.ListenAndServe()
}

func (a *RESTApiV1) Shutdown() error {
	if a.server == nil {
		return errors.New("server is not running")
	}
	return a.server.Shutdown(context.Background())
}

func (a *RESTApiV1) GetAllAPIs() []string {
	list := []string{}
	err := a.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		apiPath, err := route.GetPathTemplate()
		if err == nil {
			list = append(list, apiPath)
		}
		return err
	})
	if err != nil {
		a.logger.Error("error while getting all APIs", zap.Error(err))
	}

	return list
}
