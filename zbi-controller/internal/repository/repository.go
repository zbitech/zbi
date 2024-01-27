package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/internal/vars"
	"github.com/zbitech/controller/pkg/interfaces"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
)

var client = &http.Client{}

type RepositoryService struct {
}

func NewRepositoryService() interfaces.RepositoryServiceIF {
	return &RepositoryService{}
}

func (repo *RepositoryService) UpdateProjectResource(ctx context.Context, project string, resource *model.KubernetesResource) error {

	log := logger.GetServiceLogger(ctx, "repo.UpdateProjectResource")

	var repository = vars.ZBI_REPOSITORY_URL + "/projects/" + project + "/resources"
	jsonReq, _ := json.Marshal(resource)
	req, err := http.NewRequest(http.MethodPut, repository, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}

	log.WithFields(logrus.Fields{"repo": repository}).Infof("updating project resource")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {

		return errors.New("failed to update project resource")
	}
}

func (repo *RepositoryService) UpdateInstanceResource(ctx context.Context, instance string, resource *model.KubernetesResource) error {

	log := logger.GetServiceLogger(ctx, "repo.UpdateInstanceResource")

	//	var repository = helper.Config.GetSettings().Repository + "/projects/" + project + "/instances/" + instance + "/resources"
	var repository = vars.ZBI_REPOSITORY_URL + "/instances/" + instance + "/resources"

	jsonReq, _ := json.Marshal(resource)
	req, err := http.NewRequest(http.MethodPut, repository, bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Errorf("failed to update instance resource: %s", err)
		return err
	}

	log.WithFields(logrus.Fields{"repo": repository}).Infof("updating instance resource")

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf("failed to update instance resource")
		return errors.New("failed to update instance resource")
	}
}

func (repo *RepositoryService) GetBlockchainInfo(ctx context.Context, blockchain string) (*model.BlockchainInfo, error) {

	log := logger.GetServiceLogger(ctx, "repo.GetBlockchainInfo")
	var repository = vars.ZBI_REPOSITORY_URL + "/config/blockchains/" + blockchain

	req, err := http.NewRequest(http.MethodGet, repository, nil)
	if err != nil {
		log.Errorf("failed to get blockchain info: %s", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result model.BlockchainInfo

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, errors.New("unable to retrieve blockchain info")
		}
		return &result, nil
	} else {
		message := "failed to get blockchain info"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) GetBlockchainNodeInfo(ctx context.Context, blockchain, node string) (*model.BlockchainNodeInfo, error) {
	log := logger.GetServiceLogger(ctx, "repo.GetBlockchainInfo")
	var repository = vars.ZBI_REPOSITORY_URL + "/config/blockchains/" + blockchain + "/" + node

	req, err := http.NewRequest(http.MethodGet, repository, nil)
	if err != nil {
		log.Errorf("failed to get blockchain info: %s", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result model.BlockchainNodeInfo

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, errors.New("unable to retrieve blockchain info")
		}
		return &result, nil
	} else {
		message := "failed to get blockchain info"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) GetPolicyInfo(ctx context.Context) (*model.PolicyInfo, error) {

	log := logger.GetServiceLogger(ctx, "repo.GetPolicyInfo")
	var repository = vars.ZBI_REPOSITORY_URL + "/config/policy"

	log.Debugf("connecting to %s", repository)
	req, err := http.NewRequest(http.MethodGet, repository, nil)
	if err != nil {
		log.Errorf("failed to get policy info: %s", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result model.PolicyInfo
		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			log.Errorf("failed to read response - %s", err)
			return nil, errors.New("unable to retrieve policy info")
		}
		return &result, nil
	} else {
		message := "failed to get policy info"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) CreateProject(ctx context.Context, project *model.Project) (*model.Project, error) {

	log := logger.GetServiceLogger(ctx, "repo.GetProject")
	var repository = vars.ZBI_REPOSITORY_URL + "/projects/"

	jsonReq, _ := json.Marshal(project)
	req, err := http.NewRequest(http.MethodPost, repository, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result model.Project

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, errors.New("unable to retrieve project")
		}
		return &result, nil
	} else {
		message := "failed to get project"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}

}

func (repo *RepositoryService) GetProject(ctx context.Context, project string) (*model.Project, error) {

	log := logger.GetServiceLogger(ctx, "repo.GetProject")
	var repository = vars.ZBI_REPOSITORY_URL + "/projects/" + project

	req, err := http.NewRequest(http.MethodGet, repository, nil)
	if err != nil {
		log.Errorf("failed to get project: %s", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result model.Project

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, errors.New("unable to retrieve project")
		}
		return &result, nil
	} else {
		message := "failed to get project"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) GetProjects(ctx context.Context, owner string) ([]model.Project, error) {

	log := logger.GetServiceLogger(ctx, "repo.GetProject")
	var repository = vars.ZBI_REPOSITORY_URL + "/projects"

	req, err := http.NewRequest(http.MethodGet, repository, nil)
	if err != nil {
		log.Errorf("failed to get project: %s", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result []model.Project

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, errors.New("unable to retrieve project")
		}
		return result, nil
	} else {
		message := "failed to get project"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) GetInstance(ctx context.Context, instance string) (*model.Instance, error) {

	log := logger.GetServiceLogger(ctx, "repo.GetInstance")
	var repository = vars.ZBI_REPOSITORY_URL + "/instances/" + instance

	req, err := http.NewRequest(http.MethodGet, repository, nil)
	if err != nil {
		log.Errorf("failed to get instance: %s", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result model.Instance

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, errors.New("unable to retrieve instance")
		}
		return &result, nil
	} else {
		message := "failed to get instance"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) GetInstances(ctx context.Context, project string) ([]model.Instance, error) {

	log := logger.GetServiceLogger(ctx, "repo.GetInstance")
	var repository = vars.ZBI_REPOSITORY_URL + "/projects/" + project + "/instances"

	req, err := http.NewRequest(http.MethodGet, repository, nil)
	if err != nil {
		log.Errorf("failed to get instance: %s", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result []model.Instance

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, errors.New("unable to retrieve instance")
		}
		return result, nil
	} else {
		message := "failed to get instance"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) CreateInstance(ctx context.Context, projectId, owner string, request *model.InstanceRequest) (*model.Instance, error) {

	log := logger.GetServiceLogger(ctx, "repo.GetProject")
	var repository = vars.ZBI_REPOSITORY_URL + "/projects/" + projectId + "/instances"

	jsonReq, _ := json.Marshal(request)
	req, err := http.NewRequest(http.MethodPost, repository, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	req.Header.Add("x-owner-id", owner)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result model.Instance

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, errors.New("unable to create instance")
		}
		return &result, nil
	} else {
		message := "failed to get instance"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) UpdateInstance(ctx context.Context, instanceId string, request *model.InstanceRequest) (*model.Instance, error) {
	log := logger.GetServiceLogger(ctx, "repo.GetProject")
	var repository = vars.ZBI_REPOSITORY_URL + "/instances/" + instanceId

	jsonReq, _ := json.Marshal(request)
	req, err := http.NewRequest(http.MethodPut, repository, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result model.Instance

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, errors.New("unable to retrieve instance")
		}
		return &result, nil
	} else {
		message := "failed to get instance"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) AddProjectActivity(ctx context.Context, project string, op model.EventAction) error {

	log := logger.GetServiceLogger(ctx, "repo.AddProjectActivity")
	var repository = vars.ZBI_REPOSITORY_URL + "/projects/" + project + "/activity"

	activity := make(map[string]interface{})
	activity["op"] = op

	jsonReq, _ := json.Marshal(activity)
	req, err := http.NewRequest(http.MethodPost, repository, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		message := "failed to get instance resources"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return errors.New(message)
	}
}

func (repo *RepositoryService) AddInstanceActivity(ctx context.Context, instance string, op model.EventAction) error {

	log := logger.GetServiceLogger(ctx, "repo.AddInstanceActivity")
	var repository = vars.ZBI_REPOSITORY_URL + "/instances/" + instance + "/activity"

	activity := make(map[string]interface{})
	activity["op"] = op

	jsonReq, _ := json.Marshal(activity)
	req, err := http.NewRequest(http.MethodPost, repository, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-internal-secret", vars.ZBI_INTERNAL_CLIENT_SECRET)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		message := "failed to get instance resources"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return errors.New(message)
	}
}
