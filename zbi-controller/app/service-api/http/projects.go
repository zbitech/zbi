package http

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/app/service-api/request"
	"github.com/zbitech/controller/app/service-api/response"
	"github.com/zbitech/controller/internal/vars"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
)

func CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	var projectRequest model.Project
	if err := request.ReadJSON(w, r, &projectRequest); err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	//	projectRequest.Owner = r.Header.Get("x-owner-id") //user := ctx.Value(rctx.USERID).(string)

	repository := vars.RepositoryFactory.GetRepositoryService()
	project, err := repository.CreateProject(ctx, &projectRequest)
	if err != nil {
		log.Errorf("Failed to create project in repository %s - %s", project.Name, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()
	err = zclient.CreateProject(ctx, project)
	if err != nil {
		log.Errorf("Failed to create project %s - %s", project.Name, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddProjectActivity(ctx, project.Id, model.EventActionCreate)
	if err != nil {
		log.Errorf("failed to add create activity for project %s - %s", project.Id, err)
	}

	if err = response.JSON(w, http.StatusCreated, project); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectId := request.GetParameterValue(r, request.PATH_PARAM, "project")
	if len(projectId) == 0 {
		response.BadRequestResponse(w, r, errors.New("project is required"))
		return
	}

	repository := vars.RepositoryFactory.GetRepositoryService()
	project, err := repository.GetProject(ctx, projectId)
	if err != nil {
		log.Errorf("Failed to retrieve project %s - %s", project.Name, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	instances, err := repository.GetInstances(ctx, project.Id)

	zclient := vars.KlientFactory.GetZBIClient()
	if err := zclient.DeleteProject(ctx, project, instances); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddProjectActivity(ctx, project.Id, model.EventActionDelete)
	if err != nil {
		log.Errorf("failed to add delete activity for project %s - %s", project.Id, err)
	}

	if err := response.JSON(w, http.StatusNoContent, response.Envelope{}); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

// TODO - is this allowed?
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectId := request.GetParameterValue(r, request.PATH_PARAM, "project")
	if len(projectId) == 0 {
		response.BadRequestResponse(w, r, errors.New("project is required"))
		return
	}

	var project model.Project
	if err := request.ReadJSON(w, r, &project); err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()
	err := zclient.RepairProject(ctx, &project)
	if err != nil {
		log.Errorf("Failed to update project %s - %s", project.Name, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	envelope := response.Envelope{"project": project}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func RepairProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectId := request.GetParameterValue(r, request.PATH_PARAM, "project")
	if len(projectId) == 0 {
		response.BadRequestResponse(w, r, errors.New("project is required"))
		return
	}

	repository := vars.RepositoryFactory.GetRepositoryService()
	project, err := repository.GetProject(ctx, projectId)
	if err != nil {
		log.Errorf("Failed to retrieve project %s - %s", project.Name, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()
	err = zclient.RepairProject(ctx, project)
	if err != nil {
		log.Errorf("Failed to update project %s - %s", project.Name, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddProjectActivity(ctx, project.Id, model.EventActionRepair)
	if err != nil {
		log.Errorf("failed to add repair activity for project %s - %s", project.Id, err)
	}

	if err = response.JSON(w, http.StatusOK, project); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func GetProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	log.Infof("getting projects")
	repository := vars.RepositoryFactory.GetRepositoryService()

	owner := r.Header.Get("x-owner-id")
	projects, err := repository.GetProjects(ctx, owner)

	//	zclient := vars.KlientFactory.GetZBIClient()
	//	projects, err := zclient.GetProjects(ctx)

	if err != nil {
		log.Errorf("failed to retrieve projects")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	if err = response.JSON(w, http.StatusOK, projects); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func GetProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectId := request.GetParameterValue(r, request.PATH_PARAM, "project")
	if len(projectId) == 0 {
		response.BadRequestResponse(w, r, errors.New("project is required"))
		return
	}

	log.Infof("getting project - %s", projectId)
	repository := vars.RepositoryFactory.GetRepositoryService()
	project, err := repository.GetProject(ctx, projectId)
	if err != nil {
		log.Errorf("failed to retrieve projects")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	// envelope := response.Envelope{"project": project}
	if err = response.JSON(w, http.StatusOK, project); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

// CreateInstance creates the resources associated with the instance.
// input - the instance to be created and a list of existing instances to be
// peered with the new instance.
// response - the instance and the list of resources created.
func CreateInstance(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectId := request.GetParameterValue(r, request.PATH_PARAM, "project")
	if len(projectId) == 0 {
		response.BadRequestResponse(w, r, errors.New("project is required"))
		return
	}

	repository := vars.RepositoryFactory.GetRepositoryService()
	project, err := repository.GetProject(ctx, projectId)
	if err != nil {
		log.Errorf("failed to retrieve projects")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	var instance_req model.InstanceRequest
	if err := request.ReadJSON(w, r, &instance_req); err != nil {
		log.WithFields(logrus.Fields{"error": err, "project": projectId}).Errorf("failed to read input")
		response.BadRequestResponse(w, r, err)
		return
	}
	log.WithFields(logrus.Fields{"instance": instance_req}).Infof("instance details")

	instance, err := repository.CreateInstance(ctx, project.Id, project.Owner, &instance_req)
	if err != nil {
		log.Errorf("failed to create instance")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()
	err = zclient.CreateInstance(ctx, project, instance)
	if err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddInstanceActivity(ctx, instance.Id, model.EventActionCreate)
	if err != nil {
		log.Errorf("failed to add create activity for instance %s - %s", instance.Id, err)
	}

	if err = response.JSON(w, http.StatusCreated, instance); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func GetInstances(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	projectId := request.GetParameterValue(r, request.PATH_PARAM, "project")
	if len(projectId) == 0 {
		response.BadRequestResponse(w, r, errors.New("project is required"))
		return
	}

	repository := vars.RepositoryFactory.GetRepositoryService()
	instances, err := repository.GetInstances(ctx, projectId)
	if err != nil {
		log.Errorf("failed to retrieve instances")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	if err = response.JSON(w, http.StatusOK, instances); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}
