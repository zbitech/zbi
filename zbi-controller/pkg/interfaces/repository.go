package interfaces

import (
	"context"
	"github.com/zbitech/controller/pkg/model"
)

type RepositoryServiceIF interface {
	UpdateProjectResource(ctx context.Context, project string, resource *model.KubernetesResource) error
	UpdateInstanceResource(ctx context.Context, instance string, resource *model.KubernetesResource) error
	GetBlockchainInfo(ctx context.Context, blockchain string) (*model.BlockchainInfo, error)
	GetPolicyInfo(ctx context.Context) (*model.PolicyInfo, error)
	GetProject(ctx context.Context, project string) (*model.Project, error)
	GetInstance(ctx context.Context, instance string) (*model.Instance, error)
	GetInstanceResources(ctx context.Context, instance string) (*model.KubernetesResources, error)
}

type RepositoryServiceFactoryIF interface {
	Init(ctx context.Context)
	GetRepositoryService() RepositoryServiceIF
}
