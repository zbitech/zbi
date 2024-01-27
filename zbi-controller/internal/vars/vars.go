package vars

import (
	"context"
	"strconv"

	"github.com/zbitech/controller/internal/utils"
	"github.com/zbitech/controller/pkg/interfaces"
)

var (
	CTX                        context.Context
	ASSET_PATH_DIRECTORY       = utils.GetEnv("ASSET_PATH_DIRECTORY", "tests/files/etc/zbi")
	KUBECONFIG                 = utils.GetEnv("KUBECONFIG", "cfg/kubeconfig")
	ZBI_LOG_LEVEL              = utils.GetIntEnv("ZBI_LOG_LEVEL", 0)
	USE_KUBERNETES_CONFIG, _   = strconv.ParseBool(utils.GetEnv("USEKUBERNETESCONFIG", "false"))
	CONTROLLER_METRICS, _      = strconv.ParseBool(utils.GetEnv("METRICS", "false"))
	K8S_INCLUSTER              = true
	EXPIRATION_HOURS, _        = strconv.Atoi(utils.GetEnv("EXPIRATION_HOURS", "8760"))
	ZBI_NAMESPACE              = utils.GetEnv("ZBI_NAMESPACE", "")
	ZBI_INTERNAL_CLIENT_SECRET = utils.GetEnv("ZBI_INTERNAL_CLIENT_SECRET", "zbi-internal-client")
	ZBI_REPOSITORY_URL         = utils.GetEnv("ZBI_REPOSITORY_URL", "http://localhost:4000/api")
	HOURS_IN_YEAR              = 8760

	KlientFactory     interfaces.KlientFactoryIF
	ManagerFactory    interfaces.ResourceManagerFactoryIF
	RepositoryFactory interfaces.RepositoryServiceFactoryIF
)
