package interfaces

import (
	"context"

	"github.com/zbitech/controller/pkg/model"
)

type RepositoryServiceIF interface {
	UpdateProjectResource(ctx context.Context, project string, resource *model.KubernetesResource) error
	UpdateInstanceResource(ctx context.Context, instance string, resource *model.KubernetesResource) error

	GetBlockchainInfo(ctx context.Context, blockchain string) (*model.BlockchainInfo, error)
	GetBlockchainNodeInfo(ctx context.Context, blockchain, node string) (*model.BlockchainNodeInfo, error)
	GetPolicyInfo(ctx context.Context) (*model.PolicyInfo, error)

	CreateProject(ctx context.Context, project *model.Project) (*model.Project, error)
	GetProject(ctx context.Context, project string) (*model.Project, error)
	GetProjects(ctx context.Context, owner string) ([]model.Project, error)

	CreateInstance(ctx context.Context, projectId, owner string, request *model.InstanceRequest) (*model.Instance, error)
	UpdateInstance(ctx context.Context, instanceId string, request *model.InstanceRequest) (*model.Instance, error)

	GetInstance(ctx context.Context, instance string) (*model.Instance, error)
	GetInstances(ctx context.Context, project string) ([]model.Instance, error)

	AddProjectActivity(ctx context.Context, project string, op model.EventAction) error
	AddInstanceActivity(ctx context.Context, instance string, op model.EventAction) error
}

type RepositoryServiceFactoryIF interface {
	Init(ctx context.Context)
	GetRepositoryService() RepositoryServiceIF
}
