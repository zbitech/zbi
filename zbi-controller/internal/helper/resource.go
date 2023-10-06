package helper

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/pkg/file_template"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
)

// type TemplateConfig struct {
// 	settingsConfig   *object.Settings
// 	policyConfig     *object.PolicyConfig
// 	instanceConfig   map[model.InstanceType]object.InstanceConfig
// 	instanceTemplate map[model.InstanceType]file_template.FileTemplate
// 	appTemplate      *file_template.FileTemplate
// 	projectTemplate  *file_template.FileTemplate
// }

var CONFIG_DIR = GetEnv("ZBI_CONFIG_DIRECTORY", "/etc/zbi/conf/")
var TEMPLATE_DIR = GetEnv("ZBI_TEMPLATE_DIRECTORY", "/etc/zbi/templates/")
var CONFIG_FILE = "zbi-conf.json"

var policyInfo *model.PolicyInfo
var blockchains map[string]model.BlockchainInfo = make(map[string]model.BlockchainInfo)
var instanceTemplates map[string]file_template.FileTemplate = make(map[string]file_template.FileTemplate)
var appTemplate *file_template.FileTemplate
var projectTemplate *file_template.FileTemplate

// var Config *TemplateConfig = newTemplateConfig()

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

func LoadConfig(ctx context.Context) {

	log := logger.GetLogger(ctx)

	log.Infof("loading config files from " + CONFIG_DIR)
	policyContent, err := ReadFile((CONFIG_DIR + "policies.json"))
	if err != nil {
		panic(err)
	}

	log.WithFields(logrus.Fields{"policies": string(policyContent)}).Infof("setting policy")
	var _policyInfo model.PolicyInfo
	if err = json.Unmarshal(policyContent, &_policyInfo); err != nil {
		panic(err)
	}

	SetPolicy(&_policyInfo)

	blockchainsContent, err := ReadFile((CONFIG_DIR + "blockchains.json"))
	if err != nil {
		panic(err)
	}

	log.WithFields(logrus.Fields{"blockchains": string(blockchainsContent)}).Infof("setting policy")
	var _blockchains []model.BlockchainInfo
	if err = json.Unmarshal(blockchainsContent, &_blockchains); err != nil {
		panic(err)
	}

	SetBlockchains(ctx, _blockchains, true)
}

func SetPolicy(policy *model.PolicyInfo) {
	policyInfo = policy
}

func SetBlockchains(ctx context.Context, blockchainArray []model.BlockchainInfo, loadTemplates bool) {
	for i := 0; i < len(blockchainArray); i++ {
		blockchain := blockchainArray[i]

		SetBlockchain(ctx, &blockchain, loadTemplates)
	}
}

func SetBlockchain(ctx context.Context, blockchain *model.BlockchainInfo, loadTemplates bool) {
	blockchains[blockchain.Name] = *blockchain
	if loadTemplates {
		LoadTemplates(ctx, blockchain)
	}
}

func LoadTemplates(ctx context.Context, blockchain *model.BlockchainInfo) {

	log := logger.GetLogger(ctx)

	log.Infof("creating app template from %s", TEMPLATE_DIR+"app_templates.tmpl")
	appTemplate = file_template.CreateFilePathTemplate("app", TEMPLATE_DIR+"app_templates.tmpl", file_template.NO_FUNCS)

	log.Infof("creating project template from %s", TEMPLATE_DIR+"project_templates.tmpl")
	projectTemplate = file_template.CreateFilePathTemplate("project", TEMPLATE_DIR+"project_templates.tmpl", file_template.NO_FUNCS)

	for _, node := range blockchain.Nodes {
		nodeType := node.Type
		path := TEMPLATE_DIR + nodeType + "_templates.tmpl"

		log.Infof("creating instance %s template from %s", nodeType, path)
		instanceTemplates[nodeType] = *file_template.CreateFilePathTemplate(nodeType, path, file_template.FUNCTIONS)
	}
}

func GetPolicyInfo() *model.PolicyInfo {
	return policyInfo
}

func GetBlockchainNodeInfo(iType model.InstanceType) (*model.BlockchainNodeInfo, error) {
	bc, ok := blockchains["zcash"]
	if !ok {
		return nil, errors.New("blockchain not found")
	}

	for i := 0; i < len(bc.Nodes); i++ {
		if bc.Nodes[i].Type == string(iType) {
			return &bc.Nodes[i], nil
		}
	}

	return nil, errors.New("blockchain node not found")
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

// func newTemplateConfig() *TemplateConfig {
// 	return &TemplateConfig{
// 		instanceConfig:   map[model.InstanceType]object.InstanceConfig{},
// 		instanceTemplate: map[model.InstanceType]file_template.FileTemplate{},
// 	}
// }

// func (tc *TemplateConfig) SetConfig(app object.AppConfig) {
// 	tc.settingsConfig = &app.Settings
// 	tc.policyConfig = &app.Policy
// 	for _, ic := range app.Instances {
// 		tc.instanceConfig[ic.InstanceType] = ic
// 	}
// }

// func (tc *TemplateConfig) LoadConfig(ctx context.Context) {
// 	content, err := ReadFile(CONFIG_DIR + CONFIG_FILE)
// 	if err != nil {
// 		panic(err)
// 	}

// 	log := logger.GetLogger(ctx)
// 	log.WithFields(logrus.Fields{"content": string(content)}).Infof("setting config")
// 	var appConfig object.AppConfig
// 	if err = json.Unmarshal(content, &appConfig); err != nil {
// 		panic(err)
// 	}

// 	tc.SetConfig(appConfig)

// 	settings, _ := json.Marshal(tc)
// 	log.WithFields(logrus.Fields{"content": string(settings)}).Infof("config settings")
// }

// func (tc *TemplateConfig) LoadTemplates(ctx context.Context) {

// 	log := logger.GetLogger(ctx)

// 	log.Infof("creating app template from %s", TEMPLATE_DIR+tc.settingsConfig.Templates["app"])
// 	tc.appTemplate = file_template.CreateFilePathTemplate("app", TEMPLATE_DIR+tc.settingsConfig.Templates["app"], file_template.NO_FUNCS)

// 	log.Infof("creating project template from %s", TEMPLATE_DIR+tc.settingsConfig.Templates["project"])
// 	tc.projectTemplate = file_template.CreateFilePathTemplate("project", TEMPLATE_DIR+tc.settingsConfig.Templates["project"], file_template.NO_FUNCS)

// 	for instanceType := range tc.instanceConfig {
// 		name := string(instanceType)
// 		path := TEMPLATE_DIR + tc.settingsConfig.Templates[name]

// 		log.Infof("creating instance %s template from %s", instanceType, path)
// 		tc.instanceTemplate[instanceType] = *file_template.CreateFilePathTemplate(name, path, file_template.FUNCTIONS)
// 	}

// }

// func (tc *TemplateConfig) GetSettings() *object.Settings {
// 	return tc.settingsConfig
// }

// func (tc *TemplateConfig) GetPolicyConfig() *object.PolicyConfig {
// 	return tc.policyConfig
// }

// func (tc *TemplateConfig) GetInstanceConfig(iType model.InstanceType) (*object.InstanceConfig, error) {
// 	ic, ok := tc.instanceConfig[iType]
// 	if !ok {
// 		return nil, errors.New("instance config not found for ")
// 	}

// 	return &ic, nil
// }

// func (tc *TemplateConfig) GetAppTemplate() *file_template.FileTemplate {
// 	return tc.appTemplate
// }

// func (tc *TemplateConfig) GetProjectTemplate() *file_template.FileTemplate {
// 	return tc.projectTemplate
// }

// func (tc *TemplateConfig) GetInstanceTemplate(iType model.InstanceType) (*file_template.FileTemplate, error) {
// 	it, ok := tc.instanceTemplate[iType]
// 	if !ok {
// 		return nil, errors.New("instance config not found for ")
// 	}

// 	return &it, nil
// }

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func ReadFile(path string) ([]byte, error) {
	//name := strings.Split(filepath.Base(path), ".")[0]
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}
