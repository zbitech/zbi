package manager

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/internal/vars"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type LWDInstanceResourceManager struct {
}

func NewLWDInstanceResourceManager() interfaces.InstanceResourceManagerIF {
	return &LWDInstanceResourceManager{}
}

func (L *LWDInstanceResourceManager) CreateInstanceResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "lwd.CreateInstanceResource")
	defer func() { logger.LogServiceTime(log) }()

	if peers == nil || len(peers) != 1 {
		return nil, errors.New("lightwallet instances can only be paired with one zcash")
	}

	var dataVolumeName, dataVolumeSize string

	var request = instance.Request

	dataVolumeName = fmt.Sprintf("%s-%s", instance.Name, utils.GenerateRandomString(5, true))
	dataVolumeSize = request.Volume.Size

	policy := helper.GetPolicyInfo(ctx)
	ic, err := helper.GetBlockchainNodeInfo(ctx, instance.InstanceType)
	if err != nil {
		return nil, err
	}

	fileTemplate, err := helper.GetInstanceTemplate(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	zcashInstance := getZcashInstanceHost(peers[0].Name, project.GetNamespace())
	zcashPort := getZcashInstancePort(ctx)

	lwdSpec := model.InstanceSpec{
		Name:               instance.Name,
		ServiceAccountName: policy.ServiceAccount,
		Namespace:          project.GetNamespace(),
		Labels:             helper.CreateInstanceLabels(instance),
		DomainName:         policy.DomainName,
		DomainSecret:       policy.CertificateName,
		Envoy:              helper.CreateEnvoySpec(policy, ic.GetPort(ENVOY_PORT)),
		DataVolumeName:     dataVolumeName,
		Images: map[string]string{
			LIGHT_WALLET_IMAGE: ic.GetImageRepository(LWD_IMAGE),
		},
		Ports: map[string]int32{
			GRPC: ic.GetPort(SERVICE_PORT),
			HTTP: ic.GetPort(HTTP_PORT),
		},
		Properties: map[string]interface{}{
			ZCASH_INSTANCE_NAME: request.Properties[zcashInstanceProperty],
			ZCASH_INSTANCE:      zcashInstance,
			ZCASH_PORT:          zcashPort,
			LOG_LEVEL:           request.Properties[logLevelProperty],
		},
	}

	addLWDInstance(&peers[0], instance.Name, *request)

	var specArr []string

	specArr, err = fileTemplate.ExecuteTemplates([]string{LWD_CONF, ZCASH_CONF, ENVOY_CONF, DEPLOYMENT, SERVICE}, lwdSpec)
	if err != nil {
		log.Errorf("Lightwalletd templates failed - %s", err)
		//return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	objects, err := helper.CreateYAMLObjects(specArr /*, instance.Project, instance.Name*/)
	if err != nil {
		log.Errorf("failed to generate specs for Lightwalletd - %s", err)
		// return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	storageClass := policy.StorageClass
	var volumeSpecs = []model.VolumeSpec{
		{VolumeName: dataVolumeName, StorageClass: storageClass, Namespace: project.GetNamespace(),
			VolumeDataType: string(request.Volume.Type), DataSourceType: request.Volume.Source.Type,
			SourceName: request.Volume.Source.Ref,
			Size:       dataVolumeSize, Labels: lwdSpec.Labels},
	}

	appRsc := vars.ManagerFactory.GetAppResourceManager(ctx)
	volumes, err := appRsc.CreateVolumeResource(ctx, volumeSpecs...)
	if err != nil {
		log.Errorf("Lightwalletd volume templates failed - %s", err)
		//return nil, errs.NewApplicationError(errs.ResourceGenerationError, err)
		return nil, err
	}

	objects = append(objects, volumes...)

	obj, err := L.CreateIngressResource(ctx, projIngress, project, instance, model.EventActionCreate)
	if err != nil {
		log.Errorf("Instance ingress object creation failed - %s", err)
		return nil, err
	}

	objects = append(objects, *obj)

	var resources = make([][]unstructured.Unstructured, 0)
	resources = append(resources, objects)

	return resources, nil
}

func (L *LWDInstanceResourceManager) CreateUpdateResource(ctx context.Context, project *model.Project, instance *model.Instance, peers ...model.Instance) ([][]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "lwd.CreateUpdateResource")
	defer func() { logger.LogServiceTime(log) }()

	var err error
	fileTemplate, err := helper.GetInstanceTemplate(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	var request = instance.Request

	if peers == nil || len(peers) != 1 {
		return nil, errors.New("lightwallet instances can only be paired with one zcash")
	}

	instanceSpec := model.InstanceSpec{
		Name:      instance.Name,
		Namespace: project.GetNamespace(),
		Labels:    helper.CreateInstanceLabels(instance),
		Properties: map[string]interface{}{
			ZCASH_INSTANCE_NAME: request.Properties[zcashInstanceProperty],
			ZCASH_INSTANCE:      getZcashInstanceHost(peers[0].Name, project.GetNamespace()),
			ZCASH_PORT:          getZcashInstancePort(ctx),
			LOG_LEVEL:           request.Properties[logLevelProperty],
		},
	}

	addLWDInstance(&peers[0], instance.Name, *request)

	specArr, err := fileTemplate.ExecuteTemplates([]string{LWD_CONF, ZCASH_CONF}, instanceSpec)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("lwd templates failed")
		return nil, err
	}

	objects, err := helper.CreateYAMLObjects(specArr /*, instance.Project, instance.Name*/)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("failed to generate kubernetes resources")
		return nil, err
	}

	var resources = make([][]unstructured.Unstructured, 0)
	resources = append(resources, objects)

	return resources, nil
}

func (L *LWDInstanceResourceManager) CreateIngressResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, action model.EventAction) (*unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "lwd.CreateIngressResource")
	defer func() { logger.LogServiceTime(log) }()

	policy := helper.GetPolicyInfo(ctx)
	ic, err := helper.GetBlockchainNodeInfo(ctx, instance.InstanceType)
	if err != nil {
		return nil, err
	}

	fileTemplate, err := helper.GetInstanceTemplate(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	lwdSpec := model.InstanceSpec{
		Name:         instance.Name,
		Namespace:    project.GetNamespace(),
		Labels:       helper.CreateInstanceLabels(instance),
		DomainName:   policy.DomainName,
		DomainSecret: policy.CertificateName,
		Envoy:        helper.CreateEnvoySpec(policy, ic.GetPort(ENVOY_PORT)),
	}

	var specObj string
	if action == model.EventActionStopInstance {
		specObj, err = fileTemplate.ExecuteTemplate(INGRESS_STOPPED, lwdSpec)
	} else {
		specObj, err = fileTemplate.ExecuteTemplate(INGRESS, lwdSpec)
	}
	if err != nil {
		log.Errorf("Lightwalletd templates failed - %s", err)
		// return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	//create new ingress resource
	object, err := helper.CreateYAMLObject(specObj)
	if err != nil {
		log.Errorf("Lightwalletd templates - %s", err)
		// return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	return object, nil
}

func (L *LWDInstanceResourceManager) CreateStartResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "lwd.CreateStartResource")
	defer func() { logger.LogServiceTime(log) }()

	pvc := instance.Resources.Persistentvolumeclaim // GetResourceByType(model.ResourcePersistentVolumeClaim)

	policy := helper.GetPolicyInfo(ctx)
	ic, err := helper.GetBlockchainNodeInfo(ctx, instance.InstanceType)
	if err != nil {
		return nil, err
	}

	fileTemplate, err := helper.GetInstanceTemplate(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	lwdSpec := model.InstanceSpec{
		Name:               instance.Name,
		ServiceAccountName: policy.ServiceAccount,
		Namespace:          project.GetNamespace(),
		Labels:             helper.CreateInstanceLabels(instance),
		DomainName:         policy.DomainName,
		DomainSecret:       policy.CertificateName,
		Envoy:              helper.CreateEnvoySpec(policy, ic.GetPort(ENVOY_PORT)),
		DataVolumeName:     pvc.Name,
		Images: map[string]string{
			LIGHT_WALLET_IMAGE: ic.GetImageRepository(LWD_IMAGE),
		},
		Ports: map[string]int32{
			GRPC: ic.GetPort(SERVICE_PORT),
			HTTP: ic.GetPort(HTTP_PORT),
		},
	}

	var specArr []string
	specArr, err = fileTemplate.ExecuteTemplates([]string{"DEPLOYMENT", "SERVICE"}, lwdSpec)
	if err != nil {
		log.Errorf("Lightwalletd templates failed - %s", err)
		//return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	objects, err := helper.CreateYAMLObjects(specArr /*, instance.Project, instance.Name*/)
	if err != nil {
		log.Errorf("Lightwalletd templates failed - %s", err)
		// return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	obj, err := L.CreateIngressResource(ctx, projIngress, project, instance, model.EventActionStartInstance)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "instance": instance}).Errorf("instance ingress object creation failed")
		return nil, err
	}

	objects = append(objects, *obj)

	return objects, nil
}

func (L *LWDInstanceResourceManager) CreateStopResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]model.KubernetesResource, []unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "lwd.CreateStopResource")
	defer func() { logger.LogServiceTime(log) }()

	var resources = make([]model.KubernetesResource, 0)
	var objects = make([]unstructured.Unstructured, 0)

	deployment := instance.Resources.Deployment // GetResourceByType(model.ResourceDeployment)
	if deployment != nil && deployment.Status == "active" {
		resources = append(resources, *deployment)
	}

	if len(resources) == 0 {
		return nil, nil, errors.New("instance is not active")
	}

	obj, err := L.CreateIngressResource(ctx, projIngress, project, instance, model.EventActionStopInstance)
	if err != nil {
		return nil, nil, err
	}
	objects = append(objects, *obj)

	return resources, objects, nil
}

func (L *LWDInstanceResourceManager) CreateRepairResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance, peers ...model.Instance) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "lwd.CreateRepairResource")
	defer func() { logger.LogServiceTime(log) }()

	if peers == nil || len(peers) != 1 {
		return nil, errors.New("lightwallet instances can only be paired with one zcash")
	}

	var request = instance.Request

	var dataVolumeName, dataVolumeSize string

	pvc := instance.Resources.Persistentvolumeclaim // GetResourceByType(model.ResourcePersistentVolumeClaim)
	if pvc != nil && pvc.Status == "active" {
		dataVolumeName = pvc.Name
		dataVolumeSize = request.Volume.Size
	} else {
		dataVolumeName = fmt.Sprintf("%s-%s", instance.Name, utils.GenerateRandomString(5, true))
		dataVolumeSize = request.Volume.Size
		pvc.Name = dataVolumeName // create new volume
	}

	policy := helper.GetPolicyInfo(ctx)
	ic, err := helper.GetBlockchainNodeInfo(ctx, instance.InstanceType)
	if err != nil {
		return nil, err
	}

	fileTemplate, err := helper.GetInstanceTemplate(instance.InstanceType)
	if err != nil {
		return nil, err
	}

	zcashInstance := getZcashInstanceHost(peers[0].Name, project.GetNamespace())
	zcashPort := getZcashInstancePort(ctx)

	lwdSpec := model.InstanceSpec{
		Name:               instance.Name,
		ServiceAccountName: policy.ServiceAccount,
		Namespace:          project.GetNamespace(),
		Labels:             helper.CreateInstanceLabels(instance),
		DomainName:         policy.DomainName,
		DomainSecret:       policy.CertificateName,
		Envoy:              helper.CreateEnvoySpec(policy, ic.GetPort(ENVOY_PORT)),
		DataVolumeName:     dataVolumeName,
		Images: map[string]string{
			LIGHT_WALLET_IMAGE: ic.GetImageRepository(LWD_IMAGE),
		},
		Ports: map[string]int32{
			GRPC: ic.GetPort(SERVICE_PORT),
			HTTP: ic.GetPort(HTTP_PORT),
		},
		Properties: map[string]interface{}{
			ZCASH_INSTANCE_NAME: request.Properties[zcashInstanceProperty],
			ZCASH_INSTANCE:      zcashInstance,
			ZCASH_PORT:          zcashPort,
			LOG_LEVEL:           request.Properties[logLevelProperty],
		},
	}

	addLWDInstance(&peers[0], instance.Name, *request)

	var specArr []string

	specArr, err = fileTemplate.ExecuteTemplates([]string{LWD_CONF, ZCASH_CONF, ENVOY_CONF, DEPLOYMENT, SERVICE}, lwdSpec)
	if err != nil {
		log.Errorf("Lightwalletd templates failed - %s", err)
		//return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	objects, err := helper.CreateYAMLObjects(specArr /*, instance.Project, instance.Name*/)
	if err != nil {
		log.Errorf("failed to generate specs for Lightwalletd - %s", err)
		// return nil, errs.NewApplicationError(errs.ResourceRetrievalError, err)
		return nil, err
	}

	if pvc == nil || pvc.Status != "active" {
		storageClass := policy.StorageClass
		var volumeSpecs = []model.VolumeSpec{
			{VolumeName: dataVolumeName, StorageClass: storageClass, Namespace: project.GetNamespace(),
				VolumeDataType: string(request.Volume.Type), DataSourceType: request.Volume.Source.Type,
				SourceName: request.Volume.Source.Ref,
				Size:       dataVolumeSize, Labels: lwdSpec.Labels},
		}

		appRsc := vars.ManagerFactory.GetAppResourceManager(ctx)
		volumes, err := appRsc.CreateVolumeResource(ctx, volumeSpecs...)
		if err != nil {
			log.Errorf("Lightwalletd volume templates failed - %s", err)
			//return nil, errs.NewApplicationError(errs.ResourceGenerationError, err)
			return nil, err
		}

		objects = append(objects, volumes...)
	}

	obj, err := L.CreateIngressResource(ctx, projIngress, project, instance, model.EventActionCreate)
	if err != nil {
		log.Errorf("Instance ingress object creation failed - %s", err)
		return nil, err
	}

	objects = append(objects, *obj)

	return objects, nil
}

func (L *LWDInstanceResourceManager) CreateSnapshotResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "lwd.CreateSnapshotResource")
	defer func() { logger.LogServiceTime(log) }()

	resource := instance.Resources.Persistentvolumeclaim // GetResourceByType(model.ResourcePersistentVolumeClaim)

	policy := helper.GetPolicyInfo(ctx)

	var req model.SnapshotRequest
	req.Namespace = project.GetNamespace()
	req.VolumeName = resource.Name
	req.SnapshotClass = policy.SnapshotClass
	req.Labels = helper.CreateInstanceLabels(instance)

	appRsc := vars.ManagerFactory.GetAppResourceManager(ctx)
	return appRsc.CreateSnapshotResource(ctx, &req)
}

func (L *LWDInstanceResourceManager) CreateSnapshotScheduleResource(ctx context.Context, project *model.Project, instance *model.Instance, scheduleType model.SnapshotScheduleType) ([]unstructured.Unstructured, error) {

	var log = logger.GetServiceLogger(ctx, "lwd.CreateSnapshotScheduleResource")
	defer func() { logger.LogServiceTime(log) }()

	var req model.SnapshotScheduleRequest

	policy := helper.GetPolicyInfo(ctx)
	resource := instance.Resources.Persistentvolumeclaim // GetResourceByType(model.ResourcePersistentVolumeClaim)

	req.Namespace = project.GetNamespace()
	req.VolumeName = resource.Name
	req.Schedule = scheduleType
	req.SnapshotClass = policy.SnapshotClass
	req.Labels = helper.CreateInstanceLabels(instance)

	appRsc := vars.ManagerFactory.GetAppResourceManager(ctx)
	return appRsc.CreateSnapshotScheduleResource(ctx, &req)
}

func (L *LWDInstanceResourceManager) CreateRotationResource(ctx context.Context, project *model.Project, instance *model.Instance) ([]unstructured.Unstructured, error) {
	var log = logger.GetServiceLogger(ctx, "lwd.CreateRotationResource")
	defer func() { logger.LogServiceTime(log) }()

	// TODO - rotate certificate ?
	return []unstructured.Unstructured{}, nil
}

func (L *LWDInstanceResourceManager) CreateDeleteResource(ctx context.Context, projIngress *unstructured.Unstructured, project *model.Project, instance *model.Instance) ([]model.KubernetesResource, []unstructured.Unstructured, error) {
	var log = logger.GetServiceLogger(ctx, "lwd.CreateDeleteResource")
	defer func() { logger.LogServiceTime(log) }()

	resources := instance.GetResourceArray()
	return resources, []unstructured.Unstructured{}, nil
}
