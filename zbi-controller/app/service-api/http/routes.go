package http

import (
	"context"
	"net/http"

	"github.com/zbitech/controller/app/service-api/http/middleware"
	"github.com/zbitech/controller/app/service-api/response"
	"github.com/zbitech/controller/app/service-api/server"
	"github.com/zbitech/controller/pkg/logger"
)

func SetupRoutes(ctx context.Context, server *server.HttpServer) {

	log := logger.GetLogger(ctx)

	router := server.GetRouter()

	log.Infof("initializing middlewares")
	router.Use(middleware.InitRequest, middleware.Logging)

	router.NotFoundHandler = http.HandlerFunc(response.NotFoundResponse)
	router.MethodNotAllowedHandler = http.HandlerFunc(response.MethodNotAllowedResponse)

	log.Infof("setting project routers")
	project := router.PathPrefix("/api/projects").Subrouter()
	project.Handle("", middleware.Chain(GetProjects)).Methods(http.MethodGet)
	project.Handle("", middleware.Chain(CreateProject)).Methods(http.MethodPost)
	project.Handle("/{project}", middleware.Chain(GetProject)).Methods(http.MethodGet)
	project.Handle("/{project}", middleware.Chain(DeleteProject)).Methods(http.MethodDelete)
	project.Handle("/{project}", middleware.Chain(UpdateProject)).Methods(http.MethodPut)          // update
	project.Handle("/{project}/repair", middleware.Chain(RepairProject)).Methods(http.MethodPatch) // repair

	project.Handle("/{project}/instances", middleware.Chain(GetInstances)).Methods(http.MethodGet)
	project.Handle("/{project}/instances", middleware.Chain(CreateInstance)).Methods(http.MethodPost)

	//	log.Infof("setting instance routers")
	instances := router.PathPrefix("/api/instances").Subrouter()
	instances.Handle("/{instance}", middleware.Chain(GetInstance)).Methods(http.MethodGet)
	instances.Handle("/{instance}", middleware.Chain(UpdateInstance)).Methods(http.MethodPut)
	instances.Handle("/{instance}", middleware.Chain(DeleteInstance)).Methods(http.MethodDelete)

	instances.Handle("/{instance}/repair", middleware.Chain(RepairInstance)).Methods(http.MethodPatch)                                    // repair
	instances.Handle("/{instance}/{action:stop|start|snapshot|schedule|rotate}", middleware.Chain(PatchInstance)).Methods(http.MethodPut) // activate, deactivate, snapshot, backup

}
