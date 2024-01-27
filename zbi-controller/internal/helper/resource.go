package helper

import (
	"context"
	"encoding/base64"
	"errors"

	ttlcache "github.com/jellydator/ttlcache/v3"
	"github.com/zbitech/controller/internal/vars"
	"github.com/zbitech/controller/pkg/file_template"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
)

//var CONFIG_DIR = GetEnv("ZBI_CONFIG_DIRECTORY", "/etc/zbi/conf/")
//var TEMPLATE_DIR = GetEnv("ZBI_TEMPLATE_DIRECTORY", "/etc/zbi/templates/")
//var CONFIG_FILE = "zbi-conf.json"

var policyInfo *model.PolicyInfo
var blockchains map[string]model.BlockchainInfo = make(map[string]model.BlockchainInfo)

var blockchainnodes map[string]model.BlockchainNodeInfo = make(map[string]model.BlockchainNodeInfo)
var instanceTemplates map[string]file_template.FileTemplate = make(map[string]file_template.FileTemplate)
var appTemplate *file_template.FileTemplate
var projectTemplate *file_template.FileTemplate

var (
	NAMESPACE       = "NAMESPACE"
	ZCASH_CONF      = "ZCASH_CONF"
	ENVOY_CONF      = "ENVOY_CONF"
	LWD_CONF        = "LWD_CONF"
	CREDENTIALS     = "CREDENTIALS"
	DEPLOYMENT      = "DEPLOYMENT"
	SERVICE         = "SERVICE"
	INGRESS         = "INGRESS"
	INGRESS_INCLUDE = "INGRESS_INCLUDE"
	VOLUME          = "VOLUME"
	SNAPSHOT        = "SNAPSHOT"
	VOLUME_SNAPSHOT = "VOLUME_SNAPSHOT"
	INSTANCE_LIST   = "INSTANCE_LIST"
)

var cache *ttlcache.Cache[string, interface{}]

func CachePolicyInfo(ctx context.Context) {

	log := logger.GetLogger(ctx)
	repository := vars.RepositoryFactory.GetRepositoryService()

	log.Info("Retrieving policy information")
	policyInfo, err := repository.GetPolicyInfo(ctx)
	if err != nil {
		log.Errorf("Failed to get policy - %s", err)
	} else {
		log.Info("Caching policy information")
		cache.Set("policy", policyInfo, ttlcache.DefaultTTL)
	}

}

func CacheBlockchainInfo(ctx context.Context) {

	log := logger.GetLogger(ctx)
	repository := vars.RepositoryFactory.GetRepositoryService()

	log.Info("Retrieving zcash blockchain information")
	blockchainInfo, err := repository.GetBlockchainInfo(ctx, "zcash")
	if err != nil {
		log.Errorf("Failed to get zcash blockchain - %s", err)
	} else {

		log.Info("Caching zcash blockchain information")
		for node, content := range blockchainInfo.Templates {
			log.Infof("Creating kubernetes templates for %s", node)
			template, err := base64.StdEncoding.DecodeString(content)

			//			log.Info(string(template))
			if err != nil {
				log.Errorf("Failed to decode %s template - %s", node, err)
			} else {
				if node == "app" {
					appTemplate, _ = file_template.NewTextTemplate("app", string(template), file_template.NO_FUNCS)
				} else if node == "project" {
					projectTemplate, _ = file_template.NewTextTemplate("project", string(template), file_template.NO_FUNCS)
				} else {
					t, _ := file_template.NewTextTemplate(node, string(template), file_template.FUNCTIONS)
					instanceTemplates[node] = *t
				}
			}
		}

		for _, blockchain := range blockchainInfo.Nodes {
			blockchainnodes[blockchain.Type] = blockchain
		}

		cache.Set("blockchain", blockchainInfo, ttlcache.DefaultTTL)
	}
}

func LoadConfig(ctx context.Context) {

	cache = ttlcache.New[string, interface{}]()

	go cache.Start()

	CachePolicyInfo(ctx)
	CacheBlockchainInfo(ctx)
}

func GetPolicyInfo(ctx context.Context) *model.PolicyInfo {

	policyInfo := cache.Get("policy")
	return policyInfo.Value().(*model.PolicyInfo)
}

func GetBlockchainNodeInfo(ctx context.Context, iType model.InstanceType) (*model.BlockchainNodeInfo, error) {

	bc, ok := blockchainnodes[string(iType)]
	if !ok {
		return nil, errors.New("blockchain node not found")
	}

	return &bc, nil
}

func GetAppTemplate() *file_template.FileTemplate {
	return appTemplate
}

func GetProjectTemplate() *file_template.FileTemplate {
	return projectTemplate
}

func GetInstanceTemplate(iType model.InstanceType) (*file_template.FileTemplate, error) {
	it, ok := instanceTemplates[string(iType)]
	if !ok {
		return nil, errors.New("blockchain node template not found")
	}

	return &it, nil
}
