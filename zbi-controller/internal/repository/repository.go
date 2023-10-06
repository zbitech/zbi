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

	var repository = vars.CONTROL_PLANE_URL + "/projects/" + project + "/resources"
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
	var repository = vars.CONTROL_PLANE_URL + "/instances/" + instance + "/resources"

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
	var repository = vars.CONTROL_PLANE_URL + "/config/blockchains/" + blockchain

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
		type Response struct {
			Blockchain model.BlockchainInfo `json:"blockchain"`
		}
		var result Response

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return &result.Blockchain, nil
		}
		return nil, errors.New("unable to retrieve blockchain info")
	} else {
		message := "failed to get blockchain info"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) GetPolicyInfo(ctx context.Context) (*model.PolicyInfo, error) {

	log := logger.GetServiceLogger(ctx, "repo.GetPolicyInfo")
	var repository = vars.CONTROL_PLANE_URL + "/config/policy"

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
		type Response struct {
			Policy model.PolicyInfo `json:"policy"`
		}
		var result Response

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return &result.Policy, nil
		}
		return nil, errors.New("unable to retrieve policy info")
	} else {
		message := "failed to get policy info"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) GetProject(ctx context.Context, project string) (*model.Project, error) {

	log := logger.GetServiceLogger(ctx, "repo.GetProject")
	var repository = vars.CONTROL_PLANE_URL + "/projects/" + project

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
		type Response struct {
			Project model.Project `json:"project"`
		}
		var result Response

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return &result.Project, nil
		}
		return nil, errors.New("unable to retrieve project")
	} else {
		message := "failed to get project"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) GetInstance(ctx context.Context, instance string) (*model.Instance, error) {

	log := logger.GetServiceLogger(ctx, "repo.GetInstance")
	var repository = vars.CONTROL_PLANE_URL + "/instances/" + instance

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
		type Response struct {
			Instance model.Instance `json:"instance"`
		}
		var result Response

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return &result.Instance, nil
		}
		return nil, errors.New("unable to retrieve instance")
	} else {
		message := "failed to get instance"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}

func (repo *RepositoryService) GetInstanceResources(ctx context.Context, instance string) (*model.KubernetesResources, error) {

	log := logger.GetServiceLogger(ctx, "repo.GetInstanceResources")
	var repository = vars.CONTROL_PLANE_URL + "/instances/" + instance + "/resources"

	req, err := http.NewRequest(http.MethodGet, repository, nil)
	if err != nil {
		log.Errorf("failed to get instance resource: %s", err)
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
		type Response struct {
			Resources model.KubernetesResources `json:"resources"`
		}
		var result Response

		body, err := io.ReadAll(resp.Body)
		if err = json.Unmarshal(body, &result); err != nil {
			return &result.Resources, nil
		}
		return nil, errors.New("unable to retrieve instance resources")
	} else {
		message := "failed to get instance resources"
		log.WithFields(logrus.Fields{"status": resp.StatusCode, "detail": resp.Body}).Errorf(message)
		return nil, errors.New(message)
	}
}
