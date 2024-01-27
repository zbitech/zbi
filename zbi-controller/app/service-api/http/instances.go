package http

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/app/service-api/request"
	"github.com/zbitech/controller/app/service-api/response"
	"github.com/zbitech/controller/internal/vars"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
)

func DeleteInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	instanceId := request.GetParameterValue(r, request.PATH_PARAM, "instance")
	if len(instanceId) == 0 {
		response.BadRequestResponse(w, r, errors.New("instance is required"))
		return
	}

	repository := vars.RepositoryFactory.GetRepositoryService()
	instance, err := repository.GetInstance(ctx, instanceId)
	if err != nil {
		log.Errorf("failed to retrieve instance")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	repository.GetProject(ctx, "")
	repository.GetInstance(ctx, "")

	zclient := vars.KlientFactory.GetZBIClient()

	err = zclient.DeleteInstance(ctx, instance.Project, instance)
	if err != nil {
		log.Errorf("failed to delete instance resources - %s", err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddInstanceActivity(ctx, instance.Id, model.EventActionDelete)
	if err != nil {
		log.Errorf("failed to add snapshot activity for instance %s - %s", instance.Id, err)
	}

	if err = response.JSON(w, http.StatusOK, instance); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func UpdateInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	instanceId := request.GetParameterValue(r, request.PATH_PARAM, "instance")
	if len(instanceId) == 0 {
		response.BadRequestResponse(w, r, errors.New("instance is required"))
		return
	}

	repository := vars.RepositoryFactory.GetRepositoryService()
	instance, err := repository.GetInstance(ctx, instanceId)
	if err != nil {
		log.Errorf("failed to retrieve instance")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	var instance_req model.InstanceRequest
	if err := request.ReadJSON(w, r, &instance_req); err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instanceId}).Errorf("failed to read input")
		response.BadRequestResponse(w, r, err)
		return
	}
	log.WithFields(logrus.Fields{"instance": instance_req}).Infof("instance details")

	instance, err = repository.UpdateInstance(ctx, instance.Id, &instance_req)
	if err != nil {
		log.Errorf("failed to create instance")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	log.WithFields(logrus.Fields{"instance": instance}).Infof("updating instance")
	zclient := vars.KlientFactory.GetZBIClient()
	err = zclient.UpdateInstance(ctx, instance.Project, instance)
	if err != nil {
		//		service.HandleError(ctx, w, r, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddInstanceActivity(ctx, instance.Id, model.EventActionUpdate)
	if err != nil {
		log.Errorf("failed to add update activity for instance %s - %s", instance.Id, err)
	}

	if err = response.JSON(w, http.StatusNoContent, response.Envelope{}); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func RepairInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	instanceId := request.GetParameterValue(r, request.PATH_PARAM, "instance")
	if len(instanceId) == 0 {
		response.BadRequestResponse(w, r, errors.New("instance is required"))
		return
	}

	repository := vars.RepositoryFactory.GetRepositoryService()
	instance, err := repository.GetInstance(ctx, instanceId)
	if err != nil {
		log.Errorf("failed to retrieve instance")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	log.WithFields(logrus.Fields{"instance": instance}).Infof("repairing instance")

	zclient := vars.KlientFactory.GetZBIClient()
	err = zclient.RepairInstance(ctx, instance.Project, instance)
	if err != nil {
		//		service.HandleError(ctx, w, r, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddInstanceActivity(ctx, instance.Id, model.EventActionRepair)
	if err != nil {
		log.Errorf("failed to add repair activity for instance %s - %s", instance.Id, err)
	}

	if err = response.JSON(w, http.StatusOK, instance); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}

}

func StartInstance(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := logger.GetLogger(ctx)

	instanceId := request.GetParameterValue(r, request.PATH_PARAM, "instance")
	if len(instanceId) == 0 {
		response.BadRequestResponse(w, r, errors.New("instance is required"))
		return
	}

	repository := vars.RepositoryFactory.GetRepositoryService()
	instance, err := repository.GetInstance(ctx, instanceId)
	if err != nil {
		log.Errorf("failed to retrieve instance")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()

	err = zclient.StartInstance(ctx, instance.Project, instance)

	if err != nil {
		log.Errorf("failed to start instance %s - %s", instanceId, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddInstanceActivity(ctx, instance.Id, model.EventActionStartInstance)
	if err != nil {
		log.Errorf("failed to add start activity for instance %s - %s", instance.Id, err)
	}

	if err = response.JSON(w, http.StatusOK, instance); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func StopInstance(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := logger.GetLogger(ctx)

	instanceId := request.GetParameterValue(r, request.PATH_PARAM, "instance")
	if len(instanceId) == 0 {
		response.BadRequestResponse(w, r, errors.New("instance is required"))
		return
	}

	repository := vars.RepositoryFactory.GetRepositoryService()
	instance, err := repository.GetInstance(ctx, instanceId)
	if err != nil {
		log.Errorf("failed to retrieve instance")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()

	err = zclient.StopInstance(ctx, instance.Project, instance)
	if err != nil {
		log.Errorf("failed to stop instance %s - %s", instanceId, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddInstanceActivity(ctx, instance.Id, model.EventActionStopInstance)
	if err != nil {
		log.Errorf("failed to add stop activity for instance %s - %s", instance.Id, err)
	}

	if err = response.JSON(w, http.StatusOK, instance); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func CreateSnapshot(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := logger.GetLogger(ctx)

	instanceId := request.GetParameterValue(r, request.PATH_PARAM, "instance")
	if len(instanceId) == 0 {
		response.BadRequestResponse(w, r, errors.New("instance is required"))
		return
	}

	repository := vars.RepositoryFactory.GetRepositoryService()
	instance, err := repository.GetInstance(ctx, instanceId)
	if err != nil {
		log.Errorf("failed to retrieve instance")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()

	err = zclient.CreateSnapshot(ctx, instance.Project, instance)
	if err != nil {
		log.Errorf("failed to stop instance %s - %s", instanceId, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddInstanceActivity(ctx, instance.Id, model.EventActionSnapshot)
	if err != nil {
		log.Errorf("failed to add snapshot activity for instance %s - %s", instance.Id, err)
	}

	envelope := response.Envelope{"project": instance.Project, "instance": instance}
	if err = response.JSON(w, http.StatusOK, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func RotateInstanceCredentials(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := logger.GetLogger(ctx)

	instanceId := request.GetParameterValue(r, request.PATH_PARAM, "instance")
	if len(instanceId) == 0 {
		response.BadRequestResponse(w, r, errors.New("instance is required"))
		return
	}

	repository := vars.RepositoryFactory.GetRepositoryService()
	instance, err := repository.GetInstance(ctx, instanceId)
	if err != nil {
		log.Errorf("failed to retrieve instance")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()

	err = zclient.RotateInstanceCredentials(ctx, instance.Project, instance)
	if err != nil {
		log.Errorf("failed to stop instance %s - %s", instanceId, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddInstanceActivity(ctx, instance.Id, model.EventActionRotate)
	if err != nil {
		log.Errorf("failed to add rotate activity for instance %s - %s", instance.Id, err)
	}

	if err = response.JSON(w, http.StatusOK, instance); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func CreateSnapshotSchedule(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := logger.GetLogger(ctx)

	instanceId := request.GetParameterValue(r, request.PATH_PARAM, "instance")
	if len(instanceId) == 0 {
		response.BadRequestResponse(w, r, errors.New("instance is required"))
		return
	}

	var schedule model.SnapshotScheduleType

	if err := request.ReadJSON(w, r, &schedule); err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instanceId, "schedule": schedule}).Errorf("failed to read input")
		response.BadRequestResponse(w, r, err)
		return
	}

	repository := vars.RepositoryFactory.GetRepositoryService()
	instance, err := repository.GetInstance(ctx, instanceId)
	if err != nil {
		log.Errorf("failed to retrieve instance")
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	zclient := vars.KlientFactory.GetZBIClient()

	err = zclient.CreateSnapshotSchedule(ctx, instance.Project, instance, schedule)
	if err != nil {
		log.Errorf("failed to stop instance %s - %s", instanceId, err)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddInstanceActivity(ctx, instance.Id, model.EventActionSchedule)
	if err != nil {
		log.Errorf("failed to add snapshot schedule activity for instance %s - %s", instance.Id, err)
	}

	if err = response.JSON(w, http.StatusOK, instance); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func PatchInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	action := request.GetParameterValue(r, request.PATH_PARAM, "action")

	switch action {
	case "snapshot":
		CreateSnapshot(w, r)
		return
	case "schedule":
		CreateSnapshotSchedule(w, r)
		return
	case "start":
		StartInstance(w, r)
		return
	case "stop":
		StopInstance(w, r)
		return
	case "rotate":
		RotateInstanceCredentials(w, r)
		return
	}

	envelope := response.Envelope{"message": "invalid request"}
	if err := response.JSON(w, http.StatusBadRequest, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func GetInstance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	instanceId := request.GetParameterValue(r, request.PATH_PARAM, "instance")

	log.Infof("getting instance %s", instanceId)

	repository := vars.RepositoryFactory.GetRepositoryService()
	instance, err := repository.GetInstance(ctx, instanceId)
	if err != nil {
		log.Errorf("failed to retrieve instance %s", instanceId)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	if err = response.JSON(w, http.StatusOK, instance); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}

func DeleteInstanceResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	instanceId := request.GetParameterValue(r, request.PATH_PARAM, "instance")
	if len(instanceId) == 0 {
		response.BadRequestResponse(w, r, errors.New("instance is required"))
		return
	}

	resourceName := request.GetParameterValue(r, request.GET_PARAM, "name")
	if len(resourceName) == 0 {
		response.BadRequestResponse(w, r, errors.New("resource name is required"))
		return
	}

	resourceType := request.GetParameterValue(r, request.GET_PARAM, "type")
	if len(instanceId) == 0 {
		response.BadRequestResponse(w, r, errors.New("resource type is required"))
		return
	}

	log.Infof("getting instance %s", instanceId)

	repository := vars.RepositoryFactory.GetRepositoryService()

	instance, err := repository.GetInstance(ctx, instanceId)
	if err != nil {
		log.Errorf("failed to retrieve instance %s", instanceId)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	log.Infof("deleting resource %s (%s) for instance %s", resourceName, resourceType, instanceId)
	zclient := vars.KlientFactory.GetZBIClient()
	err = zclient.DeleteInstanceResource(ctx, instance.Project, instance, resourceName, model.ResourceObjectType(resourceType))
	if err != nil {
		log.Errorf("failed to delete resource %s (%s) for instance %s", resourceName, resourceType, instanceId)
		response.ServerErrorResponse(w, r, ctx, err)
		return
	}

	err = repository.AddInstanceActivity(ctx, instance.Id, model.EventActionDeleteResource)
	if err != nil {
		log.Errorf("failed to add delete resource (%s.%s) activity for instance %s - %s", resourceType, resourceName, instance.Id, err)
	}

	envelope := response.Envelope{}
	if err = response.JSON(w, http.StatusNoContent, envelope); err != nil {
		response.ServerErrorResponse(w, r, ctx, err)
	}
}
