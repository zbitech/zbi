package model

// func (project *Project) GetInstanceType(name string) InstanceType {
// 	for _, entry := range project.Instances {
// 		if entry.Name == name {
// 			return entry.InstanceType
// 		}
// 	}

// 	return ""
// }

func appendResource(resources []KubernetesResource, newResource KubernetesResource) []KubernetesResource {

	for index := 0; index < len(resources); index++ {
		if resources[index].Name == newResource.Name && resources[index].Type == newResource.Type {
			resources[index].Properties = newResource.Properties
			resources[index].Updated = newResource.Updated
			resources[index].Status = newResource.Status
			return resources
		}
	}

	return append(resources, newResource)
}

func (instance *Instance) GetResourceArray() []KubernetesResource {
	return CreateResourceArray(instance.Resources)
}

func CreateResourceArray(resources *KubernetesResources) []KubernetesResource {
	var array []KubernetesResource

	if resources.Configmap != nil {
		array = append(array, *resources.Configmap)
	}

	if resources.Deployment != nil {
		array = append(array, *resources.Deployment)
	}

	if resources.Httpproxy != nil {
		array = append(array, *resources.Httpproxy)
	}

	if resources.Namespace != nil {
		array = append(array, *resources.Namespace)
	}

	if resources.Persistentvolumeclaim != nil {
		array = append(array, *resources.Persistentvolumeclaim)
	}

	if resources.Secret != nil {
		array = append(array, *resources.Secret)
	}

	if resources.Service != nil {
		array = append(array, *resources.Service)
	}

	if resources.Snapshotschedule != nil {
		array = append(array, *resources.Snapshotschedule)
	}

	if resources.Volumesnapshot != nil {
		array = append(array, resources.Volumesnapshot...)
	}

	if resources.Httpproxy != nil {
		array = append(array, *resources.Httpproxy)
	}

	return array
}

// func (instance *Instance) GetResource(resourceName string, resourceType ResourceObjectType) *KubernetesResource {

// 	var resources []KubernetesResource
// 	if resourceType == ResourceVolumeSnapshot {
// 		resources = instance.Resources.Snapshots
// 	} else if resourceType == ResourceSnapshotSchedule {
// 		resources = instance.Resources.Schedules
// 	} else {
// 		resources = instance.Resources.Resources
// 	}

// 	for _, resource := range resources {
// 		if resource.Name == resourceName && resource.Type == resourceType {
// 			return &resource
// 		}
// 	}

// 	return nil
// }

// func (instance *Instance) GetResourceByName(resourceName string) *KubernetesResource {

// 	for _, resource := range instance.Resources.Resources {
// 		if resource.Name == resourceName {
// 			return &resource
// 		}
// 	}

// 	for _, resource := range instance.Resources.Snapshots {
// 		if resource.Name == resourceName {
// 			return &resource
// 		}
// 	}

// 	for _, resource := range instance.Resources.Schedules {
// 		if resource.Name == resourceName {
// 			return &resource
// 		}
// 	}

// 	return nil
// }

// func (instance *Instance) GetResourceByType(resourceType ResourceObjectType) *KubernetesResource {

// 	var resources []KubernetesResource
// 	if resourceType == ResourceVolumeSnapshot {
// 		resources = instance.Resources.Snapshots
// 	} else if resourceType == ResourceSnapshotSchedule {
// 		resources = instance.Resources.Schedules
// 	} else {
// 		resources = instance.Resources.Resources
// 	}

// 	for _, resource := range resources {
// 		if resource.Type == resourceType {
// 			return &resource
// 		}
// 	}

// 	return nil
// }

// func (instance *Instance) AddResource(resource KubernetesResource) {

// 	if resource.Type == ResourceVolumeSnapshot {
// 		instance.Resources.Snapshots = append(instance.Resources.Snapshots, resource)
// 	} else if resource.Type == ResourceSnapshotSchedule {
// 		instance.Resources.Schedules = append(instance.Resources.Schedules, resource)
// 	} else {
// 		instance.Resources.Resources = append(instance.Resources.Resources, resource)
// 	}
// }

// func (instance *Instance) AddResources(resources ...KubernetesResource) {
// 	for _, resource := range resources {
// 		instance.AddResource(resource)
// 	}
// }

// func (instance *Instance) HasResources() bool {
// 	return instance.Resources != nil && (len(instance.Resources.Resources) > 0 ||
// 		len(instance.Resources.Snapshots) > 0 ||
// 		len(instance.Resources.Schedules) > 0)
// }

// func (instance *Instance) GetResources() []KubernetesResource {
// 	var resources = make([]KubernetesResource, 0)

// 	resources = append(resources, instance.Resources.Resources...)
// 	resources = append(resources, instance.Resources.Snapshots...)
// 	resources = append(resources, instance.Resources.Schedules...)

// 	return resources
// }

func (project *Project) GetNamespace() string {
	return project.Name
}

func (project *Project) GetResourceArray() []KubernetesResource {
	return CreateResourceArray(project.Resources)
}

// func (project *Project) GetResource(resourceName string, resourceType ResourceObjectType) *KubernetesResource {
// 	for _, resource := range project.Resources {
// 		if resource.Name == resourceName && resource.Type == resourceType {
// 			return &resource
// 		}
// 	}

// 	return nil
// }

// func (project *Project) GetResourceByName(resourceName string) *KubernetesResource {
// 	for _, resource := range project.Resources {
// 		if resource.Name == resourceName {
// 			return &resource
// 		}
// 	}

// 	return nil
// }

// func (project *Project) GetResourceByType(resourceType ResourceObjectType) *KubernetesResource {
// 	for _, resource := range project.Resources {
// 		if resource.Type == resourceType {
// 			return &resource
// 		}
// 	}

// 	return nil
// }

// func (project *Project) AddResource(resource KubernetesResource) {
// 	if project.Resources == nil {
// 		project.Resources = make([]KubernetesResource, 0)
// 	}
// 	project.Resources = appendResource(project.Resources, resource)
// }

// func (project *Project) AddResources(resources ...KubernetesResource) {
// 	for _, resource := range resources {
// 		project.AddResource(resource)
// 	}
// }

func (node *BlockchainNodeInfo) GetImage(name string) *ImageInfo {
	for _, image := range node.Images {
		if image.Name == name {
			return &image
		}
	}

	return nil
}

func (node *BlockchainNodeInfo) GetImageRepository(name string) string {
	for _, image := range node.Images {
		if image.Name == name {
			return image.Url
		}
	}

	return ""
}

func (node *BlockchainNodeInfo) GetPort(name string) int32 {
	port, ok := node.Ports[name]
	if !ok {
		return -1
	}

	return port
}

// func ResourceObjectTypeToString(rtype ResourceObjectType) string {

// }

// func StringToResourceObjectType(rtype string) ResourceObjectType {

// }
